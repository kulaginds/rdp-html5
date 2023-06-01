package gcc

import (
	"errors"
	"io"

	"github.com/kulaginds/rdp-html5/internal/pkg/rdp/per"
)

type ConferenceCreateResponse struct{}

func (r *ConferenceCreateResponse) Deserialize(wire io.Reader) error {
	_, err := per.ReadChoice(wire)
	if err != nil {
		return err
	}

	var objectIdentifier bool

	objectIdentifier, err = per.ReadObjectIdentifier(t124_02_98_oid, wire)
	if err != nil {
		return err
	}

	if !objectIdentifier {
		return errors.New("bad object identifier t124")
	}

	_, err = per.ReadLength(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadChoice(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadInteger16(1001, wire)
	if err != nil {
		return err
	}

	_, err = per.ReadInteger(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadEnumerates(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadNumberOfSet(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadChoice(wire)
	if err != nil {
		return err
	}

	var octetStream bool

	octetStream, err = per.ReadOctetStream([]byte(h221SCKey), 4, wire)
	if err != nil {
		return err
	}

	if !octetStream {
		return errors.New("bad H221 SC_KEY")
	}

	_, err = per.ReadLength(wire)
	if err != nil {
		return err
	}

	return nil
}
