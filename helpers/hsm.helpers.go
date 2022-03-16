package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/miekg/pkcs11/p11"
	"gitlab.finema.co/finema/etda/key-repository-api/consts"
	"gitlab.finema.co/finema/etda/key-repository-api/emsgs"
	core "ssi-gitlab.teda.th/ssi/core"
)

func NewHSMSession(slotNumber int, pin string) (p11.Session, core.IError) {
	module, err := p11.OpenModule("/usr/local/lib/libcs_pkcs11_R2.so")
	if err != nil {
		return nil, emsgs.HSMInitializeError(err)
	}

	slots, err := module.Slots()
	if err != nil {
		return nil, emsgs.HSMSlotError(err)
	}

	if slotNumber > len(slots)-1 {
		err := errors.New("slot not found")
		return nil, emsgs.HSMSlotError(err)
	}

	slot := slots[slotNumber]

	// always close all session before open a new one
	err = slot.CloseAllSessions()
	if err != nil {
		return nil, emsgs.HSMSessionError(err)
	}

	session, err := slot.OpenSession()
	if err != nil {
		return nil, emsgs.HSMSessionError(err)
	}

	err = session.Login(pin)
	if err != nil {
		return nil, emsgs.HSMLoginError(err)
	}

	return session, nil
}

func KeepHSMAlive(contextOptions *core.ContextOptions, slotNumber int, pin string) {
	for {
		log.Println("reconnecting hsm to keep alive")
		newSession, err := NewHSMSession(slotNumber, pin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "HSM: %v", err)
			os.Exit(1)
		}
		contextOptions.DATA[consts.ContextKeyHSMSession] = newSession
		time.Sleep(12 * time.Minute)
	}
}
