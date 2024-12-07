package hivego

import "time"

func getTestVoteOp() HiveOperation {
	return voteOperation{
		Voter:    "xeroc",
		Author:   "xeroc",
		Permlink: "piston",
		Weight:   10000,
		opText:   "vote",
	}
}

func getTestCustomJsonOp() HiveOperation {
	return customJsonOperation{
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{"xeroc"},
		Id:                   "test-id",
		Json:                 "{\"testk\":\"testv\"}",
		opText:               "custom_json",
	}
}

func getTestAccountUpdateOp() HiveOperation {
	return accountUpdateOperation{
		Account:      "sniperduel17",
		Owner:        nil,
		Active:       nil,
		Posting:      nil,
		MemoKey:      "STM6n4WcwyiC63udKYR8jDFuzG9T48dhy2Qb5sVmQ9MyNuKM7xE29",
		JsonMetadata: "{\"foo\":\"bar\"}",
		opText:       "account_update",
	}
}

func getTwoTestOps() []HiveOperation {
	return []HiveOperation{getTestVoteOp(), getTestCustomJsonOp()}
}

func getTestTx(ops []HiveOperation) hiveTransaction {
	exp, _ := time.Parse("2006-01-02T15:04:05", "2016-08-08T12:24:17")
	expStr := exp.Format("2006-01-02T15:04:05")

	return hiveTransaction{
		RefBlockNum:    36029,
		RefBlockPrefix: 1164960351,
		Expiration:     expStr,
		Operations:     ops,
	}
}

func getTestVoteTx() hiveTransaction {
	return getTestTx([]HiveOperation{getTestVoteOp()})
}
