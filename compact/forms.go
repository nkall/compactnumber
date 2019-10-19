package compact

var compactFormsByLanguage = map[string]CompactForms{
	"en": map[CompactType][]CompactFormRule{
		Short: {
			{
				Type:            1000,
				ZeroesInPattern: 1,
				PatternsByPluralForm: map[string]string{
					"one":   "0K",
					"other": "0K",
				},
			},
			{
				Type:            10000,
				ZeroesInPattern: 2,
				PatternsByPluralForm: map[string]string{
					"one":   "00K",
					"other": "00K",
				},
			},
			{
				Type:            100000,
				ZeroesInPattern: 3,
				PatternsByPluralForm: map[string]string{
					"one":   "000K",
					"other": "000K",
				},
			},
			{
				Type:            1000000,
				ZeroesInPattern: 1,
				PatternsByPluralForm: map[string]string{
					"one":   "0M",
					"other": "0M",
				},
			},
			{
				Type:            10000000,
				ZeroesInPattern: 2,
				PatternsByPluralForm: map[string]string{
					"one":   "00M",
					"other": "00M",
				},
			},
			{
				Type:            100000000,
				ZeroesInPattern: 3,
				PatternsByPluralForm: map[string]string{
					"one":   "000M",
					"other": "000M",
				},
			},
			{
				Type:            1000000000,
				ZeroesInPattern: 1,
				PatternsByPluralForm: map[string]string{
					"one":   "0B",
					"other": "0B",
				},
			},
			{
				Type:            10000000000,
				ZeroesInPattern: 2,
				PatternsByPluralForm: map[string]string{
					"one":   "00B",
					"other": "00B",
				},
			},
			{
				Type:            100000000000,
				ZeroesInPattern: 3,
				PatternsByPluralForm: map[string]string{
					"one":   "000B",
					"other": "000B",
				},
			},
			{
				Type:            1000000000000,
				ZeroesInPattern: 1,
				PatternsByPluralForm: map[string]string{
					"one":   "0T",
					"other": "0T",
				},
			},
			{
				Type:            10000000000000,
				ZeroesInPattern: 2,
				PatternsByPluralForm: map[string]string{
					"one":   "00T",
					"other": "00T",
				},
			},
			{
				Type:            100000000000000,
				ZeroesInPattern: 3,
				PatternsByPluralForm: map[string]string{
					"one":   "000T",
					"other": "000T",
				},
			},
		},
	},
}
