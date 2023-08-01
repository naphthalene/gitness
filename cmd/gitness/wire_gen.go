// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/harness/gitness/cli/server"
	"github.com/harness/gitness/events"
	"github.com/harness/gitness/gitrpc"
	server3 "github.com/harness/gitness/gitrpc/server"
	"github.com/harness/gitness/gitrpc/server/cron"
	check2 "github.com/harness/gitness/internal/api/controller/check"
	"github.com/harness/gitness/internal/api/controller/githook"
	"github.com/harness/gitness/internal/api/controller/principal"
	"github.com/harness/gitness/internal/api/controller/pullreq"
	"github.com/harness/gitness/internal/api/controller/repo"
	"github.com/harness/gitness/internal/api/controller/service"
	"github.com/harness/gitness/internal/api/controller/serviceaccount"
	"github.com/harness/gitness/internal/api/controller/space"
	"github.com/harness/gitness/internal/api/controller/system"
	"github.com/harness/gitness/internal/api/controller/user"
	webhook2 "github.com/harness/gitness/internal/api/controller/webhook"
	"github.com/harness/gitness/internal/auth/authn"
	"github.com/harness/gitness/internal/auth/authz"
	"github.com/harness/gitness/internal/bootstrap"
	events3 "github.com/harness/gitness/internal/events/git"
	events2 "github.com/harness/gitness/internal/events/pullreq"
	"github.com/harness/gitness/internal/router"
	server2 "github.com/harness/gitness/internal/server"
	"github.com/harness/gitness/internal/services"
	"github.com/harness/gitness/internal/services/codecomments"
	pullreq2 "github.com/harness/gitness/internal/services/pullreq"
	"github.com/harness/gitness/internal/services/webhook"
	"github.com/harness/gitness/internal/store"
	"github.com/harness/gitness/internal/store/cache"
	"github.com/harness/gitness/internal/store/database"
	"github.com/harness/gitness/internal/url"
	"github.com/harness/gitness/lock"
	"github.com/harness/gitness/pubsub"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/check"
)

// Injectors from wire.go:

func initSystem(ctx context.Context, config *types.Config) (*server.System, error) {
	principalUID := check.ProvidePrincipalUIDCheck()
	databaseConfig := server.ProvideDatabaseConfig(config)
	db, err := database.ProvideDatabase(ctx, databaseConfig)
	if err != nil {
		return nil, err
	}
	pathTransformation := store.ProvidePathTransformation()
	pathStore := database.ProvidePathStore(db, pathTransformation)
	pathCache := cache.ProvidePathCache(pathStore, pathTransformation)
	spaceStore := database.ProvideSpaceStore(db, pathCache)
	principalInfoView := database.ProvidePrincipalInfoView(db)
	principalInfoCache := cache.ProvidePrincipalInfoCache(principalInfoView)
	membershipStore := database.ProvideMembershipStore(db, principalInfoCache)
	permissionCache := authz.ProvidePermissionCache(spaceStore, membershipStore)
	authorizer := authz.ProvideAuthorizer(permissionCache)
	principalUIDTransformation := store.ProvidePrincipalUIDTransformation()
	principalStore := database.ProvidePrincipalStore(db, principalUIDTransformation)
	tokenStore := database.ProvideTokenStore(db)
	controller := user.NewController(principalUID, authorizer, principalStore, tokenStore, membershipStore)
	serviceController := service.NewController(principalUID, authorizer, principalStore)
	bootstrapBootstrap := bootstrap.ProvideBootstrap(config, controller, serviceController)
	authenticator := authn.ProvideAuthenticator(principalStore, tokenStore)
	provider, err := url.ProvideURLProvider(config)
	if err != nil {
		return nil, err
	}
	pathUID := check.ProvidePathUIDCheck()
	repoStore := database.ProvideRepoStore(db, pathCache)
	gitrpcConfig, err := server.ProvideGitRPCClientConfig()
	if err != nil {
		return nil, err
	}
	gitrpcInterface, err := gitrpc.ProvideClient(gitrpcConfig)
	if err != nil {
		return nil, err
	}
	repoController := repo.ProvideController(config, db, provider, pathUID, authorizer, pathStore, repoStore, spaceStore, principalStore, gitrpcInterface)
	spaceController := space.ProvideController(db, provider, pathUID, authorizer, pathStore, spaceStore, repoStore, principalStore, repoController, membershipStore)
	pullReqStore := database.ProvidePullReqStore(db, principalInfoCache)
	pullReqActivityStore := database.ProvidePullReqActivityStore(db, principalInfoCache)
	codeCommentView := database.ProvideCodeCommentView(db)
	pullReqReviewStore := database.ProvidePullReqReviewStore(db)
	pullReqReviewerStore := database.ProvidePullReqReviewerStore(db, principalInfoCache)
	eventsConfig, err := server.ProvideEventsConfig()
	if err != nil {
		return nil, err
	}
	universalClient, err := server.ProvideRedis(config)
	if err != nil {
		return nil, err
	}
	eventsSystem, err := events.ProvideSystem(eventsConfig, universalClient)
	if err != nil {
		return nil, err
	}
	reporter, err := events2.ProvideReporter(eventsSystem)
	if err != nil {
		return nil, err
	}
	lockConfig := server.ProvideLockConfig(config)
	mutexManager := lock.ProvideMutexManager(lockConfig, universalClient)
	migrator := codecomments.ProvideMigrator(gitrpcInterface)
	pullreqController := pullreq.ProvideController(db, provider, authorizer, pullReqStore, pullReqActivityStore, codeCommentView, pullReqReviewStore, pullReqReviewerStore, repoStore, principalStore, gitrpcInterface, reporter, mutexManager, migrator)
	webhookConfig, err := server.ProvideWebhookConfig()
	if err != nil {
		return nil, err
	}
	webhookStore := database.ProvideWebhookStore(db)
	webhookExecutionStore := database.ProvideWebhookExecutionStore(db)
	readerFactory, err := events3.ProvideReaderFactory(eventsSystem)
	if err != nil {
		return nil, err
	}
	eventsReaderFactory, err := events2.ProvideReaderFactory(eventsSystem)
	if err != nil {
		return nil, err
	}
	webhookService, err := webhook.ProvideService(ctx, webhookConfig, readerFactory, eventsReaderFactory, webhookStore, webhookExecutionStore, repoStore, pullReqStore, provider, principalStore, gitrpcInterface)
	if err != nil {
		return nil, err
	}
	webhookController := webhook2.ProvideController(webhookConfig, db, authorizer, webhookStore, webhookExecutionStore, repoStore, webhookService)
	eventsReporter, err := events3.ProvideReporter(eventsSystem)
	if err != nil {
		return nil, err
	}
	githookController := githook.ProvideController(db, authorizer, principalStore, repoStore, eventsReporter)
	serviceaccountController := serviceaccount.NewController(principalUID, authorizer, principalStore, spaceStore, repoStore, tokenStore)
	principalController := principal.ProvideController(principalStore)
	checkStore := database.ProvideCheckStore(db, principalInfoCache)
	checkController := check2.ProvideController(db, authorizer, repoStore, checkStore, gitrpcInterface)
	systemController := system.NewController(principalStore, config)
	apiHandler := router.ProvideAPIHandler(config, authenticator, repoController, spaceController, pullreqController, webhookController, githookController, serviceaccountController, controller, principalController, checkController, systemController)
	gitHandler := router.ProvideGitHandler(config, provider, repoStore, authenticator, authorizer, gitrpcInterface)
	webHandler := router.ProvideWebHandler(config)
	routerRouter := router.ProvideRouter(config, apiHandler, gitHandler, webHandler)
	serverServer := server2.ProvideServer(config, routerRouter)
	serverConfig, err := server.ProvideGitRPCServerConfig()
	if err != nil {
		return nil, err
	}
	gitAdapter, err := server3.ProvideGITAdapter()
	if err != nil {
		return nil, err
	}
	grpcServer, err := server3.ProvideServer(serverConfig, gitAdapter)
	if err != nil {
		return nil, err
	}
	manager := cron.ProvideManager(serverConfig)
	repoGitInfoView := database.ProvideRepoGitInfoView(db)
	repoGitInfoCache := cache.ProvideRepoGitInfoCache(repoGitInfoView)
	pubsubConfig := pubsub.ProvideConfig(config)
	pubSub := pubsub.ProvidePubSub(pubsubConfig, universalClient)
	pullreqService, err := pullreq2.ProvideService(ctx, config, readerFactory, eventsReaderFactory, reporter, gitrpcInterface, db, repoGitInfoCache, repoStore, pullReqStore, pullReqActivityStore, codeCommentView, migrator, pubSub, provider)
	if err != nil {
		return nil, err
	}
	servicesServices := services.ProvideServices(webhookService, pullreqService)
	serverSystem := server.NewSystem(bootstrapBootstrap, serverServer, grpcServer, manager, servicesServices)
	return serverSystem, nil
}
