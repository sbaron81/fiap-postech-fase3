package main

import "testing"

func TestGenerateAPIKeyHasExpectedPrefixAndLength(t *testing.T) {
	key, err := generateAPIKey()
	if err != nil {
		t.Fatalf("generateAPIKey retornou erro: %v", err)
	}

	const prefix = "tm_key_"
	if len(key) <= len(prefix) || key[:len(prefix)] != prefix {
		t.Errorf("esperava prefixo %q, recebeu %q", prefix, key)
	}

	// 32 bytes em hex = 64 caracteres, mais o prefixo.
	wantLen := len(prefix) + 64
	if len(key) != wantLen {
		t.Errorf("esperava chave com %d caracteres, recebeu %d (%q)", wantLen, len(key), key)
	}
}

func TestGenerateAPIKeyIsRandom(t *testing.T) {
	key1, err := generateAPIKey()
	if err != nil {
		t.Fatalf("generateAPIKey retornou erro: %v", err)
	}
	key2, err := generateAPIKey()
	if err != nil {
		t.Fatalf("generateAPIKey retornou erro: %v", err)
	}

	if key1 == key2 {
		t.Errorf("duas chamadas geraram a mesma chave: %q", key1)
	}
}

func TestHashAPIKeyIsDeterministic(t *testing.T) {
	hash1 := hashAPIKey("minha-chave-de-teste")
	hash2 := hashAPIKey("minha-chave-de-teste")

	if hash1 != hash2 {
		t.Errorf("o mesmo input gerou hashes diferentes: %q != %q", hash1, hash2)
	}
	if len(hash1) != 64 {
		t.Errorf("esperava hash SHA-256 em hex (64 chars), recebeu %d chars", len(hash1))
	}
}

func TestHashAPIKeyDiffersForDifferentInputs(t *testing.T) {
	hash1 := hashAPIKey("chave-a")
	hash2 := hashAPIKey("chave-b")

	if hash1 == hash2 {
		t.Errorf("inputs diferentes geraram o mesmo hash: %q", hash1)
	}
}
