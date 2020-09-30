package spec

// DESADV message specification.
var DESADV = Msg("DESADV",
	S("UNA", C, 1),
	S("UNB", C, 1),
	S("UNH", M, 1),
	S("BGM", M, 1),
	S("DTM", C, 10),
	S("ALI", C, 5),
	S("MEA", C, 5),
	S("MOA", C, 5),
	SG("SG1", C, 10,
		S("RFF", M, 1),
		S("DTM", C, 1),
	),
	SG("SG2", C, 99,
		S("NAD", M, 1),
		S("LOC", C, 10),
		SG("SG3", C, 10,
			S("RFF", M, 1),
		),
		SG("SG4", C, 10,
			S("CTA", M, 1),
			S("COM", C, 5),
		),
	),
	SG("SG5", C, 10,
		S("TOD", M, 1),
		S("LOC", C, 5),
	),
	SG("SG6", C, 10,
		S("TDT", M, 1),
		SG("SG7", C, 10,
			S("LOG", M, 1),
			S("DTM", C, 10),
		),
	),
	SG("SG8", C, 10,
		S("EQD", M, 1),
		S("MEA", C, 5),
		S("SEL", C, 5),
	),
	SG("SG10", C, 9999,
		S("CPS", M, 1),
		S("FTX", C, 5),
		SG("SG11", C, 9999,
			S("PAK", M, 1),
			S("MEA", C, 10),
			S("QTY", C, 10),
			SG("SG12", C, 10,
				S("HAN", M, 1),
			),
			SG("SG13", C, 1000,
				S("PCI", M, 1),
				S("RFF", C, 1),
				S("DTM", C, 5),
				SG("SG15", C, 99,
					S("GIN", M, 1),
				),
			),
		),
		SG("SG17", C, 9999,
			S("LIN", M, 1),
			S("PIA", C, 10),
			S("IMD", C, 25),
			S("MEA", C, 10),
			S("QTY", C, 10),
			S("ALI", C, 10),
			S("DLM", C, 100),
			S("DTM", C, 5),
			S("FTX", C, 99),
			S("MOA", C, 5),
			SG("SG18", C, 99,
				S("RFF", M, 1),
				S("DTM", C, 1),
			),
			SG("SG20", C, 100,
				S("LOC", M, 1),
				S("NAD", C, 1),
				S("DTM", C, 1),
				S("QTY", C, 10),
			),
			SG("SG22", C, 9999,
				S("PCI", M, 1),
				S("DTM", C, 5),
				S("MEA", C, 10),
				S("QTY", C, 1),
				SG("SG23", C, 10,
					S("GIN", M, 1),
					S("DLM", C, 100),
				),
				SG("SG24", C, 10,
					S("HAN", M, 1),
				),
			),
			SG("SG25", C, 10,
				S("QVR", M, 1),
				S("DTM", C, 5),
			),
		),
	),
	S("CNT", C, 5),
	S("UNT", M, 1),
)
