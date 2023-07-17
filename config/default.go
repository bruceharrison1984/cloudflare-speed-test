package config

import "github.com/bruceharrison1984/cloudflare-speed-test/types"

/* Retreive default test cases */
func GetDefaultConfig() ([]types.SpeedTestCase, int) {
	downloadTestCases := make([]types.SpeedTestCase, 10)
	downloadTestCases[0] = types.SpeedTestCase{PayloadSize: 1e5, Iterations: 10, TestType: types.Download}
	downloadTestCases[1] = types.SpeedTestCase{PayloadSize: 1e6, Iterations: 8, TestType: types.Download}
	downloadTestCases[2] = types.SpeedTestCase{PayloadSize: 1e7, Iterations: 6, TestType: types.Download}
	downloadTestCases[2] = types.SpeedTestCase{PayloadSize: 2.5e7, Iterations: 6, TestType: types.Download}
	downloadTestCases[3] = types.SpeedTestCase{PayloadSize: 1e8, Iterations: 3, TestType: types.Download}

	downloadTestCases[4] = types.SpeedTestCase{PayloadSize: 1e5, Iterations: 8, TestType: types.Upload}
	downloadTestCases[5] = types.SpeedTestCase{PayloadSize: 1e6, Iterations: 6, TestType: types.Upload}
	downloadTestCases[6] = types.SpeedTestCase{PayloadSize: 1e7, Iterations: 4, TestType: types.Upload}
	downloadTestCases[7] = types.SpeedTestCase{PayloadSize: 2.5e7, Iterations: 4, TestType: types.Upload}
	downloadTestCases[8] = types.SpeedTestCase{PayloadSize: 5e7, Iterations: 3, TestType: types.Upload}

	var iterations int
	for i := 0; i < len(downloadTestCases); i++ {
		iterations += downloadTestCases[i].Iterations
	}

	return downloadTestCases, iterations
}
