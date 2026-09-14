package main

import "testing"

func TestRunEvaluationLogicFlagDisabledReturnsFalse(t *testing.T) {
	a := &App{}
	info := &CombinedFlagInfo{
		Flag: &Flag{Name: "my-flag", IsEnabled: false},
	}

	if got := a.runEvaluationLogic(info, "user-1"); got != false {
		t.Errorf("flag desabilitada deveria retornar false, recebeu %v", got)
	}
}

func TestRunEvaluationLogicNoFlagReturnsFalse(t *testing.T) {
	a := &App{}
	info := &CombinedFlagInfo{Flag: nil}

	if got := a.runEvaluationLogic(info, "user-1"); got != false {
		t.Errorf("sem flag deveria retornar false, recebeu %v", got)
	}
}

func TestRunEvaluationLogicEnabledWithoutRuleReturnsTrue(t *testing.T) {
	a := &App{}
	info := &CombinedFlagInfo{
		Flag: &Flag{Name: "my-flag", IsEnabled: true},
		Rule: nil,
	}

	if got := a.runEvaluationLogic(info, "user-1"); got != true {
		t.Errorf("flag habilitada sem regra deveria retornar true, recebeu %v", got)
	}
}

func TestRunEvaluationLogicRuleDisabledReturnsTrue(t *testing.T) {
	a := &App{}
	info := &CombinedFlagInfo{
		Flag: &Flag{Name: "my-flag", IsEnabled: true},
		Rule: &TargetingRule{IsEnabled: false},
	}

	if got := a.runEvaluationLogic(info, "user-1"); got != true {
		t.Errorf("regra desabilitada deveria retornar true (so a flag manda), recebeu %v", got)
	}
}

func TestRunEvaluationLogicPercentageZeroAlwaysFalse(t *testing.T) {
	a := &App{}
	info := &CombinedFlagInfo{
		Flag: &Flag{Name: "my-flag", IsEnabled: true},
		Rule: &TargetingRule{
			IsEnabled: true,
			Rules:     Rule{Type: "PERCENTAGE", Value: float64(0)},
		},
	}

	if got := a.runEvaluationLogic(info, "qualquer-user"); got != false {
		t.Errorf("percentage 0 deveria sempre retornar false, recebeu %v para user 'qualquer-user'", got)
	}
}

func TestRunEvaluationLogicPercentageHundredAlwaysTrue(t *testing.T) {
	a := &App{}
	info := &CombinedFlagInfo{
		Flag: &Flag{Name: "my-flag", IsEnabled: true},
		Rule: &TargetingRule{
			IsEnabled: true,
			Rules:     Rule{Type: "PERCENTAGE", Value: float64(100)},
		},
	}

	if got := a.runEvaluationLogic(info, "qualquer-user"); got != true {
		t.Errorf("percentage 100 deveria sempre retornar true, recebeu %v para user 'qualquer-user'", got)
	}
}

func TestRunEvaluationLogicPercentageInvalidValueReturnsFalse(t *testing.T) {
	a := &App{}
	info := &CombinedFlagInfo{
		Flag: &Flag{Name: "my-flag", IsEnabled: true},
		Rule: &TargetingRule{
			IsEnabled: true,
			// Value nao e float64 (ex: veio como string de um JSON malformado)
			Rules: Rule{Type: "PERCENTAGE", Value: "50"},
		},
	}

	if got := a.runEvaluationLogic(info, "user-1"); got != false {
		t.Errorf("valor de percentage invalido deveria retornar false (falha segura), recebeu %v", got)
	}
}

func TestGetDeterministicBucketIsStableForSameInput(t *testing.T) {
	b1 := getDeterministicBucket("user-1my-flag")
	b2 := getDeterministicBucket("user-1my-flag")

	if b1 != b2 {
		t.Errorf("o mesmo input gerou buckets diferentes: %d != %d", b1, b2)
	}
	if b1 < 0 || b1 > 99 {
		t.Errorf("bucket fora do intervalo [0,99]: %d", b1)
	}
}

func TestGetDeterministicBucketVariesAcrossInputs(t *testing.T) {
	seen := map[int]bool{}
	for i := 0; i < 50; i++ {
		bucket := getDeterministicBucket(string(rune('a'+i)) + "my-flag")
		seen[bucket] = true
	}

	// Com 50 inputs diferentes, esperamos alguma distribuicao (nao tudo no mesmo bucket).
	if len(seen) < 5 {
		t.Errorf("distribuicao suspeita: apenas %d buckets distintos em 50 inputs", len(seen))
	}
}
