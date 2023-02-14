package server

// type Server interface {
// 	ListenandServe() error
// 	Shutdown(ctx context.Context) error
// }

// type server struct {
// 	serv *http.Server
// }

// func New() Server {
// 	conf, err := config.Load()
// 	if err != nil {
// 		log.Fatalf("failed to load config: %v", err)
// 	}

// 	r := chi.NewRouter()

// 	db, err := ConnectDB()
// 	if err != nil {
// 		log.Fatalf("failed to create database: %v", err)
// 	}

// 	// init repositories
// 	npr := repo.NewRepository(db)
// 	// guestsRepo := guests.NewRepository(db)

// 	// init services
// 	partyService := party.NewParty(npr)
// 	// guestsSrv := guests.NewService(guestsRepo, tablesSrv)

// 	// init controllers
// 	tablesCtrl := handlers.NewController(partyService)
// 	// guestsCtrl := guests.NewController(guestsHdl, guestsSrv)

// 	PR := handlers.NewPartyResource(db)
// }