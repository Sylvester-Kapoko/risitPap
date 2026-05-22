package store

import (
    "testing"
    "github.com/Sylvester-Kapoko/risitPap/domain"
)

func TestStoreConfig_RegisterID(t *testing.T) {
    st, err := Open(":memory:")
    if err != nil {
        t.Fatal(err)
    }
    defer st.db.Close()

    cfg := &domain.StoreConfig{
        StoreName:  "Test",
        StoreTaxID: "A123456789X",
    }
    // Save twice – RegisterID should only be set once
    if err := st.SaveConfig(cfg); err != nil {
        t.Fatal(err)
    }
    first := cfg.RegisterID
    if first == "" {
        t.Fatal("RegisterID was not generated")
    }

    // Save again – ID must not change
    cfg.StoreName = "New Name"
    if err := st.SaveConfig(cfg); err != nil {
        t.Fatal(err)
    }
    if cfg.RegisterID != first {
        t.Errorf("RegisterID changed from %q to %q", first, cfg.RegisterID)
    }
}