package engine_test

import (
	"encoding/json"
	"errors"
	"time"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db/dbfakes"
	"github.com/concourse/atc/db/lock/lockfakes"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/dbng/dbngfakes"
	. "github.com/concourse/atc/engine"
	"github.com/concourse/atc/engine/enginefakes"
)

var _ = Describe("DBEngine", func() {
	var (
		logger lager.Logger

		fakeEngineA *enginefakes.FakeEngine
		fakeEngineB *enginefakes.FakeEngine
		dbBuild     *dbngfakes.FakeBuild

		dbEngine Engine
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")

		fakeEngineA = new(enginefakes.FakeEngine)
		fakeEngineA.NameReturns("fake-engine-a")

		fakeEngineB = new(enginefakes.FakeEngine)
		fakeEngineB.NameReturns("fake-engine-b")

		dbBuild = new(dbngfakes.FakeBuild)
		dbBuild.IDReturns(128)

		dbEngine = NewDBEngine(Engines{fakeEngineA, fakeEngineB})
	})

	Describe("CreateBuild", func() {
		var (
			plan atc.Plan

			createdBuild Build
			buildErr     error

			planFactory atc.PlanFactory
		)

		BeforeEach(func() {
			planFactory = atc.NewPlanFactory(123)

			plan = planFactory.NewPlan(atc.TaskPlan{
				Config: &atc.TaskConfig{
					Image: "some-image",

					Params: map[string]string{
						"FOO": "1",
						"BAR": "2",
					},

					Run: atc.TaskRunConfig{
						Path: "some-script",
						Args: []string{"arg1", "arg2"},
					},
				},
			})

			dbBuild.StartReturns(true, nil)
		})

		JustBeforeEach(func() {
			createdBuild, buildErr = dbEngine.CreateBuild(logger, dbBuild, plan)
		})

		Context("when creating the build succeeds", func() {
			var fakeBuild *enginefakes.FakeBuild

			BeforeEach(func() {
				fakeBuild = new(enginefakes.FakeBuild)
				fakeBuild.MetadataReturns("some-metadata")

				fakeEngineA.CreateBuildReturns(fakeBuild, nil)
			})

			It("succeeds", func() {
				Expect(buildErr).NotTo(HaveOccurred())
			})

			It("returns a build", func() {
				Expect(createdBuild).NotTo(BeNil())
			})

			It("starts the build in the database", func() {
				Expect(dbBuild.StartCallCount()).To(Equal(1))

				engine, metadata := dbBuild.StartArgsForCall(0)
				Expect(engine).To(Equal("fake-engine-a"))
				Expect(metadata).To(Equal("some-metadata"))
			})

			Context("when the build fails to transition to started", func() {
				BeforeEach(func() {
					dbBuild.StartReturns(false, nil)
				})

				It("aborts the build", func() {
					Expect(fakeBuild.AbortCallCount()).To(Equal(1))
				})
			})
		})

		Context("when creating the build fails", func() {
			disaster := errors.New("failed")

			BeforeEach(func() {
				fakeEngineA.CreateBuildReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(buildErr).To(Equal(disaster))
			})

			It("does not start the build", func() {
				Expect(dbBuild.StartCallCount()).To(Equal(0))
			})
		})
	})

	Describe("LookupBuild", func() {
		var (
			foundBuild Build
			lookupErr  error
		)

		JustBeforeEach(func() {
			foundBuild, lookupErr = dbEngine.LookupBuild(logger, dbBuild)
		})

		It("succeeds", func() {
			Expect(lookupErr).NotTo(HaveOccurred())
		})

		It("returns a build", func() {
			Expect(foundBuild).NotTo(BeNil())
		})

		Describe("Abort", func() {
			var abortErr error

			BeforeEach(func() {
				dbBuild.ReloadReturns(true, nil)
			})

			JustBeforeEach(func() {
				abortErr = foundBuild.Abort(lagertest.NewTestLogger("test"))
			})

			Context("when acquiring the lock succeeds", func() {
				var fakeLock *lockfakes.FakeLock

				BeforeEach(func() {
					fakeLock = new(lockfakes.FakeLock)
					dbBuild.AcquireTrackingLockReturns(fakeLock, true, nil)
				})

				It("succeeds", func() {
					Expect(abortErr).NotTo(HaveOccurred())
				})

				It("marks the build as aborted", func() {
					Expect(dbBuild.AbortCallCount()).To(Equal(1))
				})
			})

			Context("when acquiring the lock fails", func() {
				BeforeEach(func() {
					dbBuild.AcquireTrackingLockReturns(nil, false, nil)
				})

				It("succeeds", func() {
					Expect(abortErr).NotTo(HaveOccurred())
				})

				It("marks the build as aborted", func() {
					Expect(dbBuild.AbortCallCount()).To(Equal(1))
				})
			})

			Context("when acquiring the lock errors", func() {
				BeforeEach(func() {
					dbBuild.AcquireTrackingLockReturns(nil, false, errors.New("bad bad bad"))
				})

				It("fails", func() {
					Expect(abortErr).To(HaveOccurred())
				})

				It("does not mark the build as aborted", func() {
					Expect(dbBuild.AbortCallCount()).To(Equal(0))
				})
			})
		})
	})

	Describe("Builds", func() {
		var build Build

		BeforeEach(func() {
			var err error
			build, err = dbEngine.LookupBuild(logger, dbBuild)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("Abort", func() {
			var abortErr error

			JustBeforeEach(func() {
				abortErr = build.Abort(lagertest.NewTestLogger("test"))
			})

			Context("when acquiring the lock succeeds", func() {
				var fakeLock *lockfakes.FakeLock

				BeforeEach(func() {
					fakeLock = new(lockfakes.FakeLock)
					dbBuild.AcquireTrackingLockReturns(fakeLock, true, nil)
				})

				Context("when the build is active", func() {
					BeforeEach(func() {
						dbBuild.ReloadReturns(true, nil)
						dbBuild.EngineReturns("fake-engine-b")

						dbBuild.AbortStub = func() error {
							Expect(dbBuild.AcquireTrackingLockCallCount()).To(Equal(1))

							_, interval := dbBuild.AcquireTrackingLockArgsForCall(0)
							Expect(interval).To(Equal(time.Minute))

							Expect(fakeLock.ReleaseCallCount()).To(BeZero())

							return nil
						}
					})

					Context("when the engine build exists", func() {
						var realBuild *enginefakes.FakeBuild

						BeforeEach(func() {
							dbBuild.ReloadReturns(true, nil)

							realBuild = new(enginefakes.FakeBuild)
							fakeEngineB.LookupBuildReturns(realBuild, nil)
						})

						Context("when aborting the db build succeeds", func() {
							BeforeEach(func() {
								dbBuild.AbortReturns(nil)
							})

							It("succeeds", func() {
								Expect(abortErr).NotTo(HaveOccurred())
							})

							It("releases the lock", func() {
								Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
							})

							It("aborts the build via the db", func() {
								Expect(dbBuild.AbortCallCount()).To(Equal(1))
							})

							It("aborts the real build", func() {
								Expect(realBuild.AbortCallCount()).To(Equal(1))
							})
						})

						Context("when aborting the db build fails", func() {
							disaster := errors.New("oh no!")

							BeforeEach(func() {
								dbBuild.AbortReturns(disaster)
							})

							It("returns the error", func() {
								Expect(abortErr).To(Equal(disaster))
							})

							It("does not abort the real build", func() {
								Expect(realBuild.AbortCallCount()).To(BeZero())
							})

							It("releases the lock", func() {
								Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
							})
						})

						Context("when aborting the real build fails", func() {
							disaster := errors.New("oh no!")

							BeforeEach(func() {
								realBuild.AbortReturns(disaster)
							})

							It("returns the error", func() {
								Expect(abortErr).To(Equal(disaster))
							})

							It("releases the lock", func() {
								Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
							})
						})
					})

					Context("when looking up the engine build fails", func() {
						disaster := errors.New("nope")

						BeforeEach(func() {
							dbBuild.ReloadReturns(true, nil)
							fakeEngineB.LookupBuildReturns(nil, disaster)
						})

						It("returns the error", func() {
							Expect(abortErr).To(Equal(disaster))
						})

						It("releases the lock", func() {
							Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
						})
					})
				})

				Context("when the build is not yet active", func() {
					BeforeEach(func() {
						dbBuild.ReloadReturns(true, nil)
						dbBuild.EngineReturns("")
					})

					It("succeeds", func() {
						Expect(abortErr).NotTo(HaveOccurred())
					})

					It("aborts the build in the db", func() {
						Expect(dbBuild.AbortCallCount()).To(Equal(1))
					})

					It("finishes the build in the db so that the aborted event is emitted", func() {
						Expect(dbBuild.FinishCallCount()).To(Equal(1))

						status := dbBuild.FinishArgsForCall(0)
						Expect(status).To(Equal(dbng.BuildStatusAborted))
					})

					It("releases the lock", func() {
						Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
					})
				})

				Context("when the build is no longer in the database", func() {
					BeforeEach(func() {
						dbBuild.ReloadReturns(false, nil)
					})

					It("succeeds", func() {
						Expect(abortErr).NotTo(HaveOccurred())
					})

					It("aborts the build in the db", func() {
						Expect(dbBuild.AbortCallCount()).To(Equal(1))
					})

					It("does not finish the build", func() {
						Expect(dbBuild.FinishCallCount()).To(Equal(0))
					})

					It("releases the lock", func() {
						Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
					})
				})
			})

			Context("when acquiring the lock errors", func() {
				BeforeEach(func() {
					dbBuild.AcquireTrackingLockReturns(nil, false, errors.New("bad bad bad"))
				})

				It("errors", func() {
					Expect(abortErr).To(HaveOccurred())
				})

				It("does not abort the build in the db", func() {
					Expect(dbBuild.AbortCallCount()).To(Equal(0))
				})
			})

			Context("when acquiring the lock fails", func() {
				BeforeEach(func() {
					dbBuild.AcquireTrackingLockReturns(nil, false, nil)
				})

				Context("when aborting the build in the db succeeds", func() {
					BeforeEach(func() {
						dbBuild.AbortReturns(nil)
					})

					It("succeeds", func() {
						Expect(abortErr).NotTo(HaveOccurred())
					})

					It("aborts the build in the db", func() {
						Expect(dbBuild.AbortCallCount()).To(Equal(1))
					})

					It("does not abort the real build", func() {
						Expect(dbBuild.ReloadCallCount()).To(BeZero())
						Expect(fakeEngineB.LookupBuildCallCount()).To(BeZero())
					})
				})

				Context("when aborting the build in the db fails", func() {
					disaster := errors.New("oh no!")

					BeforeEach(func() {
						dbBuild.AbortReturns(disaster)
					})

					It("fails", func() {
						Expect(abortErr).To(Equal(disaster))
					})
				})
			})
		})

		Describe("Resume", func() {
			var logger lager.Logger

			BeforeEach(func() {
				logger = lagertest.NewTestLogger("test")
			})

			JustBeforeEach(func() {
				build.Resume(logger)
			})

			Context("when acquiring the lock succeeds", func() {
				var fakeLock *lockfakes.FakeLock

				BeforeEach(func() {
					fakeLock = new(lockfakes.FakeLock)
					dbBuild.AcquireTrackingLockReturns(fakeLock, true, nil)
				})

				Context("when the build is active", func() {
					BeforeEach(func() {
						dbBuild.EngineReturns("fake-engine-b")
						dbBuild.IsRunningReturns(true)
						dbBuild.ReloadReturns(true, nil)
					})

					Context("when the engine build exists", func() {
						var realBuild *enginefakes.FakeBuild

						BeforeEach(func() {
							dbBuild.ReloadReturns(true, nil)

							realBuild = new(enginefakes.FakeBuild)
							fakeEngineB.LookupBuildReturns(realBuild, nil)

							realBuild.ResumeStub = func(lager.Logger) {
								Expect(dbBuild.AcquireTrackingLockCallCount()).To(Equal(1))

								_, interval := dbBuild.AcquireTrackingLockArgsForCall(0)
								Expect(interval).To(Equal(time.Minute))

								Expect(fakeLock.ReleaseCallCount()).To(BeZero())
							}
						})

						Context("when builds are released", func() {
							BeforeEach(func() {
								readyToRelease := make(chan struct{})

								go func() {
									<-readyToRelease
									dbEngine.ReleaseAll(logger)
								}()

								relased := make(chan struct{})

								realBuild.ResumeStub = func(lager.Logger) {
									close(readyToRelease)
									<-relased
								}

								fakeEngineB.ReleaseAllStub = func(lager.Logger) {
									close(relased)
								}

								aborts := make(chan struct{})
								notifier := new(dbfakes.FakeNotifier)
								notifier.NotifyReturns(aborts)

								dbBuild.AbortNotifierReturns(notifier, nil)
							})

							It("releases build engine builds", func() {
								Expect(fakeEngineB.ReleaseAllCallCount()).To(Equal(1))
							})

							It("releases the lock", func() {
								Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
							})
						})

						Context("when listening for aborts succeeds", func() {
							var (
								notifier *dbfakes.FakeNotifier
								abort    chan<- struct{}
							)

							BeforeEach(func() {
								aborts := make(chan struct{})
								abort = aborts

								notifier = new(dbfakes.FakeNotifier)
								notifier.NotifyReturns(aborts)

								dbBuild.AbortNotifierReturns(notifier, nil)
							})

							It("listens for aborts", func() {
								Expect(dbBuild.AbortNotifierCallCount()).To(Equal(1))
							})

							It("resumes the build", func() {
								Expect(realBuild.ResumeCallCount()).To(Equal(1))
							})

							It("releases the lock", func() {
								Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
							})

							It("closes the notifier", func() {
								Expect(notifier.CloseCallCount()).To(Equal(1))
							})

							Context("when the build is aborted", func() {
								var errAborted = errors.New("aborted")

								BeforeEach(func() {
									aborted := make(chan error)

									realBuild.AbortStub = func(lager.Logger) error {
										aborted <- errAborted
										return nil
									}

									realBuild.ResumeStub = func(lager.Logger) {
										<-aborted
									}

									close(abort)
								})

								It("aborts the build", func() {
									Expect(realBuild.AbortCallCount()).To(Equal(1))
								})

								It("releases the lock", func() {
									Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
								})

								It("closes the notifier", func() {
									Expect(notifier.CloseCallCount()).To(Equal(1))
								})
							})
						})

						Context("when listening for aborts fails", func() {
							disaster := errors.New("oh no!")

							BeforeEach(func() {
								dbBuild.AbortNotifierReturns(nil, disaster)
							})

							It("does not resume the build", func() {
								Expect(realBuild.ResumeCallCount()).To(BeZero())
							})

							It("releases the lock", func() {
								Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
							})
						})
					})

					Context("when looking up the engine build fails", func() {
						disaster := errors.New("nope")

						BeforeEach(func() {
							dbBuild.ReloadReturns(true, nil)
							fakeEngineB.LookupBuildReturns(nil, disaster)
						})

						It("releases the lock", func() {
							Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
						})

						It("marks the build as errored", func() {
							Expect(dbBuild.FinishCallCount()).To(Equal(1))
							buildStatus := dbBuild.FinishArgsForCall(0)
							Expect(buildStatus).To(Equal(dbng.BuildStatusErrored))
						})
					})
				})

				Context("when the build's engine is unknown", func() {
					BeforeEach(func() {
						dbBuild.ReloadReturns(true, nil)
						dbBuild.IsRunningReturns(true)
						dbBuild.EngineReturns("bogus")
					})

					It("marks the build as errored", func() {
						Expect(dbBuild.FinishCallCount()).To(Equal(1))
						buildStatus := dbBuild.FinishArgsForCall(0)
						Expect(buildStatus).To(Equal(dbng.BuildStatusErrored))
					})
				})

				Context("when the build is not yet active", func() {
					BeforeEach(func() {
						dbBuild.ReloadReturns(true, nil)
						dbBuild.EngineReturns("")
					})

					It("does not look up the build in the engine", func() {
						Expect(fakeEngineB.LookupBuildCallCount()).To(BeZero())
					})

					It("releases the lock", func() {
						Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
					})
				})

				Context("when the build has already finished", func() {
					BeforeEach(func() {
						dbBuild.ReloadReturns(true, nil)
						dbBuild.EngineReturns("fake-engine-b")
						dbBuild.StatusReturns(dbng.BuildStatusSucceeded)
					})

					It("does not look up the build in the engine", func() {
						Expect(fakeEngineB.LookupBuildCallCount()).To(BeZero())
					})

					It("releases the lock", func() {
						Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
					})
				})

				Context("when the build is no longer in the database", func() {
					BeforeEach(func() {
						dbBuild.ReloadReturns(false, nil)
					})

					It("does not look up the build in the engine", func() {
						Expect(fakeEngineB.LookupBuildCallCount()).To(BeZero())
					})

					It("releases the lock", func() {
						Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
					})
				})
			})

			Context("when acquiring the lock fails", func() {
				BeforeEach(func() {
					dbBuild.AcquireTrackingLockReturns(nil, false, errors.New("no lock for you"))
				})

				It("does not look up the build", func() {
					Expect(dbBuild.ReloadCallCount()).To(BeZero())
					Expect(fakeEngineB.LookupBuildCallCount()).To(BeZero())
				})
			})
		})

		Describe("PublicPlan", func() {
			var logger lager.Logger

			var publicPlan atc.PublicBuildPlan
			var publicPlanErr error

			BeforeEach(func() {
				logger = lagertest.NewTestLogger("test")
			})

			JustBeforeEach(func() {
				publicPlan, publicPlanErr = build.PublicPlan(logger)
			})

			var fakeLock *lockfakes.FakeLock

			BeforeEach(func() {
				fakeLock = new(lockfakes.FakeLock)
				dbBuild.AcquireTrackingLockReturns(fakeLock, true, nil)
			})

			Context("when the build is active", func() {
				BeforeEach(func() {
					dbBuild.EngineReturns("fake-engine-b")
					dbBuild.ReloadReturns(true, nil)
				})

				Context("when the engine build exists", func() {
					var realBuild *enginefakes.FakeBuild

					BeforeEach(func() {
						realBuild = new(enginefakes.FakeBuild)

						fakeEngineB.LookupBuildReturns(realBuild, nil)
					})

					Context("when getting the plan via the engine succeeds", func() {
						BeforeEach(func() {
							var plan json.RawMessage = []byte("lol")

							realBuild.PublicPlanReturns(atc.PublicBuildPlan{
								Schema: "some-schema",
								Plan:   &plan,
							}, nil)
						})

						It("succeeds", func() {
							Expect(publicPlanErr).ToNot(HaveOccurred())
						})

						It("returns the public plan from the engine", func() {
							var plan json.RawMessage = []byte("lol")

							Expect(publicPlan).To(Equal(atc.PublicBuildPlan{
								Schema: "some-schema",
								Plan:   &plan,
							}))
						})
					})

					Context("when getting the plan via the engine fails", func() {
						disaster := errors.New("nope")

						BeforeEach(func() {
							realBuild.PublicPlanReturns(atc.PublicBuildPlan{}, disaster)
						})

						It("returns the error", func() {
							Expect(publicPlanErr).To(Equal(disaster))
						})
					})
				})

				Context("when looking up the engine build fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeEngineB.LookupBuildReturns(nil, disaster)
					})

					It("returns the error", func() {
						Expect(publicPlanErr).To(Equal(disaster))
					})
				})
			})

			Context("when the build's engine is unknown", func() {
				BeforeEach(func() {
					dbBuild.EngineReturns("bogus")
				})

				It("returns an UnknownEngineError", func() {
					Expect(publicPlanErr).To(Equal(UnknownEngineError{"bogus"}))
				})
			})

			Context("when the build is not yet active", func() {
				BeforeEach(func() {
					dbBuild.ReloadReturns(true, nil)
					dbBuild.EngineReturns("")
				})

				It("does not look up the build in the engine", func() {
					Expect(fakeEngineB.LookupBuildCallCount()).To(BeZero())
				})
			})

			Context("when the build is no longer in the database", func() {
				BeforeEach(func() {
					dbBuild.ReloadReturns(false, nil)
				})

				It("does not look up the build in the engine", func() {
					Expect(fakeEngineB.LookupBuildCallCount()).To(BeZero())
				})
			})
		})
	})
})
