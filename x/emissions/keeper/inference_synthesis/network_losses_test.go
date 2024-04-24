package inference_synthesis_test

import (
	cosmosMath "cosmossdk.io/math"
	alloraMath "github.com/allora-network/allora-chain/math"

	"github.com/allora-network/allora-chain/x/emissions/keeper/inference_synthesis"
	emissions "github.com/allora-network/allora-chain/x/emissions/types"
)

func (s *InferenceSynthesisTestSuite) TestRunningWeightedAvgUpdate() {
	tests := []struct {
		name                string
		initialWeightedLoss inference_synthesis.WorkerRunningWeightedLoss
		weight              inference_synthesis.Weight
		nextValue           inference_synthesis.Weight
		epsilon             inference_synthesis.Weight
		expectedLoss        inference_synthesis.WorkerRunningWeightedLoss
		expectedErr         error
	}{
		{
			name:                "normal operation",
			initialWeightedLoss: inference_synthesis.WorkerRunningWeightedLoss{Loss: alloraMath.MustNewDecFromString("0.5"), SumWeight: alloraMath.MustNewDecFromString("1.0")},
			weight:              alloraMath.MustNewDecFromString("1.0"),
			nextValue:           alloraMath.MustNewDecFromString("2.0"),
			epsilon:             alloraMath.MustNewDecFromString("1e-4"),
			expectedLoss:        inference_synthesis.WorkerRunningWeightedLoss{Loss: alloraMath.MustNewDecFromString("0.400514997"), SumWeight: alloraMath.MustNewDecFromString("2.0")},
			expectedErr:         nil,
		},
		{
			name:                "simple example",
			initialWeightedLoss: inference_synthesis.WorkerRunningWeightedLoss{Loss: alloraMath.MustNewDecFromString("0"), SumWeight: alloraMath.MustNewDecFromString("0")},
			weight:              alloraMath.MustNewDecFromString("1.0"),
			nextValue:           alloraMath.MustNewDecFromString("0.1"),
			epsilon:             alloraMath.MustNewDecFromString("1e-4"),
			expectedLoss:        inference_synthesis.WorkerRunningWeightedLoss{Loss: alloraMath.MustNewDecFromString("-1.0"), SumWeight: alloraMath.MustNewDecFromString("1.0")},
			expectedErr:         nil,
		},
		{
			name:                "simple example2",
			initialWeightedLoss: inference_synthesis.WorkerRunningWeightedLoss{Loss: alloraMath.MustNewDecFromString("0"), SumWeight: alloraMath.MustNewDecFromString("0")},
			weight:              alloraMath.MustNewDecFromString("1.0"),
			nextValue:           alloraMath.MustNewDecFromString("0.2"),
			epsilon:             alloraMath.MustNewDecFromString("1e-4"),
			expectedLoss:        inference_synthesis.WorkerRunningWeightedLoss{Loss: alloraMath.MustNewDecFromString("-0.698970004336018804"), SumWeight: alloraMath.MustNewDecFromString("1.0")},
			expectedErr:         nil,
		},
		{
			name:                "division by zero error",
			initialWeightedLoss: inference_synthesis.WorkerRunningWeightedLoss{Loss: alloraMath.MustNewDecFromString("1.01"), SumWeight: alloraMath.MustNewDecFromString("0")},
			weight:              alloraMath.MustNewDecFromString("2.0"),
			nextValue:           alloraMath.MustNewDecFromString("1.0"),
			epsilon:             alloraMath.MustNewDecFromString("3.0"),
			expectedLoss:        inference_synthesis.WorkerRunningWeightedLoss{},
			expectedErr:         emissions.ErrFractionDivideByZero,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			updatedLoss, err := inference_synthesis.RunningWeightedAvgUpdate(
				&tc.initialWeightedLoss,
				tc.weight,
				tc.nextValue,
				tc.epsilon,
			)
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr, "Error should match the expected error")
			} else {
				s.Require().NoError(err, "No error expected but got one")
				s.Require().True(alloraMath.InDelta(tc.expectedLoss.Loss, updatedLoss.Loss, alloraMath.MustNewDecFromString("0.00001")), "Loss should match the expected value within epsilon")
				s.Require().Equal(tc.expectedLoss.SumWeight, updatedLoss.SumWeight, "Sum of weights should match the expected value")
			}
		})
	}
}

func (s *InferenceSynthesisTestSuite) TestCalcCombinedNetworkLoss() {
	tests := []struct {
		name            string
		stakesByReputer map[inference_synthesis.Worker]cosmosMath.Uint
		reportedLosses  *emissions.ReputerValueBundles
		epsilon         alloraMath.Dec
		expectedLoss    inference_synthesis.Loss
		expectedErr     error
	}{
		{
			name: "Simple case with one reputer",
			stakesByReputer: map[inference_synthesis.Worker]cosmosMath.Uint{
				"worker1": inference_synthesis.CosmosUintOneE18(), // 1 token
			},
			reportedLosses: &emissions.ReputerValueBundles{
				ReputerValueBundles: []*emissions.ReputerValueBundle{
					{
						ValueBundle: &emissions.ValueBundle{
							Reputer:       "worker1",
							CombinedValue: alloraMath.MustNewDecFromString("0.1"), // Log value of loss
						},
					},
				},
			},
			epsilon:      alloraMath.MustNewDecFromString("1e-4"),
			expectedLoss: alloraMath.MustNewDecFromString("0.1"), // exp(0.1) ≈ 1.258925
			expectedErr:  nil,
		},
		{
			name: "Two reputers",
			stakesByReputer: map[inference_synthesis.Worker]cosmosMath.Uint{
				"worker1": inference_synthesis.CosmosUintOneE18(),                            // 1 token
				"worker2": inference_synthesis.CosmosUintOneE18().Mul(cosmosMath.NewUint(2)), // 2 tokens
			},
			reportedLosses: &emissions.ReputerValueBundles{
				ReputerValueBundles: []*emissions.ReputerValueBundle{
					{
						ValueBundle: &emissions.ValueBundle{
							Reputer:       "worker1",
							CombinedValue: alloraMath.MustNewDecFromString("0.1"), // Log value of loss
						},
					},
					{
						ValueBundle: &emissions.ValueBundle{
							Reputer:       "worker2",
							CombinedValue: alloraMath.MustNewDecFromString("0.2"), // Log value of loss
						},
					},
				},
			},
			epsilon:      alloraMath.MustNewDecFromString("1e-4"),
			expectedLoss: alloraMath.MustNewDecFromString("0.1587401051968199"), // exp(0.1) ≈ 1.258925
			expectedErr:  nil,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			loss, err := inference_synthesis.CalcCombinedNetworkLoss(tc.stakesByReputer, tc.reportedLosses, tc.epsilon)

			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
			} else {
				s.Require().NoError(err)
				s.Require().True(alloraMath.InDelta(tc.expectedLoss, loss, alloraMath.MustNewDecFromString("0.00001")), "Loss should match expected value within a small epsilon")
			}
		})
	}
}

func (s *InferenceSynthesisTestSuite) TestCalcNetworkLosses() {
	tests := []struct {
		name            string
		stakesByReputer map[inference_synthesis.Worker]cosmosMath.Uint
		reportedLosses  emissions.ReputerValueBundles
		epsilon         alloraMath.Dec
		expectedOutput  emissions.ValueBundle
		expectedError   error
	}{
		{
			name: "simple two reputer combined loss",
			stakesByReputer: map[inference_synthesis.Worker]cosmosMath.Uint{
				"worker1": inference_synthesis.CosmosUintOneE18(),                            // 1 token
				"worker2": inference_synthesis.CosmosUintOneE18().Mul(cosmosMath.NewUint(2)), // 2 tokens
			},
			reportedLosses: emissions.ReputerValueBundles{
				ReputerValueBundles: []*emissions.ReputerValueBundle{
					{
						ValueBundle: &emissions.ValueBundle{
							Reputer:       "worker1",
							CombinedValue: alloraMath.MustNewDecFromString("0.1"),
							NaiveValue:    alloraMath.MustNewDecFromString("0.1"),
							InfererValues: []*emissions.WorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
							},
							ForecasterValues: []*emissions.WorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
							},
							OneOutInfererValues: []*emissions.WithheldWorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
							},
							OneOutForecasterValues: []*emissions.WithheldWorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
							},
							OneInForecasterValues: []*emissions.WorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.1"),
								},
							},
						},
					},
					{
						ValueBundle: &emissions.ValueBundle{
							Reputer:       "worker2",
							CombinedValue: alloraMath.MustNewDecFromString("0.2"),
							NaiveValue:    alloraMath.MustNewDecFromString("0.2"),
							InfererValues: []*emissions.WorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
							},
							ForecasterValues: []*emissions.WorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
							},
							OneOutInfererValues: []*emissions.WithheldWorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
							},
							OneOutForecasterValues: []*emissions.WithheldWorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
							},
							OneInForecasterValues: []*emissions.WorkerAttributedValue{
								{
									Worker: "worker1",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
								{
									Worker: "worker2",
									Value:  alloraMath.MustNewDecFromString("0.2"),
								},
							},
						},
					},
				},
			},
			epsilon: alloraMath.MustNewDecFromString("1e-4"),
			expectedOutput: emissions.ValueBundle{
				CombinedValue: alloraMath.MustNewDecFromString("0.1587401051968199"),
				NaiveValue:    alloraMath.MustNewDecFromString("0.1587401051968199"),
				InfererValues: []*emissions.WorkerAttributedValue{
					{
						Worker: "worker1",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
					{
						Worker: "worker2",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
				},
				ForecasterValues: []*emissions.WorkerAttributedValue{
					{
						Worker: "worker1",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
					{
						Worker: "worker2",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
				},
				OneOutInfererValues: []*emissions.WithheldWorkerAttributedValue{
					{
						Worker: "worker1",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
					{
						Worker: "worker2",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
				},
				OneOutForecasterValues: []*emissions.WithheldWorkerAttributedValue{
					{
						Worker: "worker1",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
					{
						Worker: "worker2",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
				},
				OneInForecasterValues: []*emissions.WorkerAttributedValue{
					{
						Worker: "worker1",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
					{
						Worker: "worker2",
						Value:  alloraMath.MustNewDecFromString("0.1587401051968199"),
					},
				},
			},
			expectedError: nil,
		},
	}

	require := s.Require()

	for _, tc := range tests {
		s.Run(tc.name, func() {
			output, err := inference_synthesis.CalcNetworkLosses(tc.stakesByReputer, tc.reportedLosses, tc.epsilon)
			if tc.expectedError != nil {
				require.Error(err)
				require.EqualError(err, tc.expectedError.Error())
			} else {
				require.NoError(err)

				require.True(alloraMath.InDelta(tc.expectedOutput.CombinedValue, output.CombinedValue, alloraMath.MustNewDecFromString("0.00001")))
				require.True(alloraMath.InDelta(tc.expectedOutput.NaiveValue, output.NaiveValue, alloraMath.MustNewDecFromString("0.00001")))

				if tc.expectedOutput.InfererValues != nil {
					require.Len(output.InfererValues, len(tc.expectedOutput.InfererValues))
					for i, expectedValue := range tc.expectedOutput.InfererValues {
						require.True(alloraMath.InDelta(expectedValue.Value, output.InfererValues[i].Value, alloraMath.MustNewDecFromString("0.00001")))
					}
				}
				if tc.expectedOutput.ForecasterValues != nil {
					require.Len(output.ForecasterValues, len(tc.expectedOutput.ForecasterValues))
					for i, expectedValue := range tc.expectedOutput.ForecasterValues {
						require.True(alloraMath.InDelta(expectedValue.Value, output.ForecasterValues[i].Value, alloraMath.MustNewDecFromString("0.00001")))
					}
				}
				if tc.expectedOutput.OneOutInfererValues != nil {
					require.Len(output.OneOutInfererValues, len(tc.expectedOutput.OneOutInfererValues))
					for i, expectedValue := range tc.expectedOutput.OneOutInfererValues {
						require.True(alloraMath.InDelta(expectedValue.Value, output.OneOutInfererValues[i].Value, alloraMath.MustNewDecFromString("0.00001")))
					}
				}
				if tc.expectedOutput.OneOutForecasterValues != nil {
					require.Len(output.OneOutForecasterValues, len(tc.expectedOutput.OneOutForecasterValues))
					for i, expectedValue := range tc.expectedOutput.OneOutForecasterValues {
						require.True(alloraMath.InDelta(expectedValue.Value, output.OneOutForecasterValues[i].Value, alloraMath.MustNewDecFromString("0.00001")))
					}
				}
				if tc.expectedOutput.OneInForecasterValues != nil {
					require.Len(output.OneInForecasterValues, len(tc.expectedOutput.OneInForecasterValues))
					for i, expectedValue := range tc.expectedOutput.OneInForecasterValues {
						require.True(alloraMath.InDelta(expectedValue.Value, output.OneInForecasterValues[i].Value, alloraMath.MustNewDecFromString("0.00001")))
					}
				}
			}
		})
	}
}
