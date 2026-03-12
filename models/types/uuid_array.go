package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// UUIDArray is a PostgreSQL uuid[] compatible type.
type UUIDArray []uuid.UUID

func (a *UUIDArray) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("UUIDArray.Scan: unsupported type")
	}

	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")
	if str == "" {
		*a = UUIDArray{}
		return nil
	}

	parts := strings.Split(str, ",")
	*a = make(UUIDArray, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(strings.Trim(s, `"`))
		if s == "" {
			continue
		}
		u, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("UUIDArray.Scan: invalid UUID %q: %w", s, err)
		}
		*a = append(*a, u)
	}
	return nil
}

func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	parts := make([]string, len(a))
	for i, u := range a {
		parts[i] = fmt.Sprintf(`"%s"`, u.String())
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (UUIDArray) GormDataType() string {
	return "uuid[]"
}
