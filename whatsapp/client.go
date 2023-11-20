package whatsapp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Client struct {
	*whatsmeow.Client
	clientLog waLog.Logger
}

func NewClient(db *sql.DB) (*Client, error) {
	store.SetOSInfo("whatsapp-openai", [3]uint32{0, 1, 0})
	dbLog := waLog.Stdout("Database", "INFO", true)
	clientLog := waLog.Stdout("Client", "INFO", true)
	storeContainer, err := newContainer(db, "sqlite3", dbLog)
	if err != nil {
		clientLog.Errorf("failed to connect to database: %v", err)
		return nil, err
	}
	deviceStore, err := storeContainer.GetFirstDevice()
	if err != nil {
		clientLog.Errorf("failed to get device: %v", err)
		return nil, err
	}
	return &Client{
		Client:    whatsmeow.NewClient(deviceStore, clientLog),
		clientLog: clientLog,
	}, nil
}

func newContainer(db *sql.DB, dialect string, log waLog.Logger) (*sqlstore.Container, error) {
	container := sqlstore.NewWithDB(db, dialect, log)
	err := container.Upgrade()
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade database: %w", err)
	}
	return container, nil
}

func (c Client) Start(ctx context.Context) error {
	if c.Store.ID == nil {
		qr, err := c.getQR(ctx)
		if err != nil {
			return err
		}
		log.Info("new qr:")
		fmt.Println(qr)
	} else {
		err := c.Connect()
		if err != nil {
			c.clientLog.Errorf("failed to connect: %v", err)
			return err
		}
	}
	return nil
}

func (c Client) getQR(ctx context.Context) (string, error) {
	ch, _ := c.GetQRChannel(ctx)
	err := c.Connect()
	if err != nil {
		c.clientLog.Errorf("failed to connect: %v", err)
		return "", err
	}
	for evt := range ch {
		if evt.Event == "code" {
			qr, err := qrcode.New(evt.Code, qrcode.Low)
			if err != nil {
				c.clientLog.Errorf("failed to encode QR: %v", err)
				return "", err
			}
			qr.DisableBorder = true
			return qr.ToSmallString(true), nil
		} else {
			log.Infof("login event: %s", evt.Event)
		}
	}
	return "", errors.New("qr code not found")
}
