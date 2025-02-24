package XPN

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	_ "github.com/xtls/xray-core/main/distro/all"

	v2net "github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf/serial"
)

type MyLogger interface {
	LogData(s string)
}

var cInstance *core.Instance

func LetsPrepareToSetMimLimit(config []byte, myLogger MyLogger) {
	myLogger.LogData("Preparing to set mem limit")
}

func SetMemLimit() {
	debug.SetGCPercent(10)
	debug.SetMemoryLimit(30 * 1024 * 1024)
}

func LetsPrepareToStart(config []byte, myLogger MyLogger) {
	myLogger.LogData("Preparing to start")
}

func LetsStart(config []byte, myLogger MyLogger) error {
	conf, err := serial.DecodeJSONConfig(bytes.NewReader(config))
	if err != nil {
		myLogger.LogData("Config load error: " + err.Error())
		return err
	}
	pbConfig, err := conf.Build()
	if err != nil {
		return err
	}
	instance, err := core.New(pbConfig)
	if err != nil {
		myLogger.LogData("Create XPN error: " + err.Error())
		return err
	}
	err = instance.Start()
	if err != nil {
		myLogger.LogData("Start XPN error: " + err.Error())
	}
	cInstance = instance
	return nil
}

func LetsPrepareToStop(myLogger MyLogger) {
	myLogger.LogData("Preparing to stop")
}

func LetsStop() {
	cInstance.Close()
}

func LetsPrepareToGetVersion(myLogger MyLogger) {
	myLogger.LogData("Preparing to get version")
}

func LetsGetVersion() string {
	return core.Version()
}

func LetsPrepareToMeasureDelay(myLogger MyLogger) {
	myLogger.LogData("Preparing to get measure delay")
}

func LetsMeasureDelay(url string) (int64, error) {
	delay, err := letsMeasureInstDelay(context.Background(), cInstance, url)
	return delay, err
}

func LetsMeasureOutboundDelay(ConfigureFileContent string, url string) (int64, error) {
	config, err := serial.LoadJSONConfig(strings.NewReader(ConfigureFileContent))
	if err != nil {
		return -1, err
	}

	// dont listen to anything for test purpose
	config.Inbound = nil
	// config.App: (fakedns), log, dispatcher, InboundConfig, OutboundConfig, (stats), router, dns, (policy)
	// keep only basic features
	config.App = config.App[:5]

	inst, err := core.New(config)
	if err != nil {
		return -1, err
	}

	inst.Start()
	delay, err := letsMeasureInstDelay(context.Background(), inst, url)
	inst.Close()
	return delay, err
}

func letsMeasureInstDelay(ctx context.Context, inst *core.Instance, url string) (int64, error) {
	if inst == nil {
		return -1, errors.New("core instance nil")
	}

	tr := &http.Transport{
		TLSHandshakeTimeout: 6 * time.Second,
		DisableKeepAlives:   true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dest, err := v2net.ParseDestination(fmt.Sprintf("%s:%s", network, addr))
			if err != nil {
				return nil, err
			}
			return core.Dial(ctx, inst, dest)
		},
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   12 * time.Second,
	}

	if len(url) <= 0 {
		url = "https://www.google.com/generate_204"
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	start := time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return -1, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return -1, fmt.Errorf("status != 20x: %s", resp.Status)
	}
	resp.Body.Close()
	return time.Since(start).Milliseconds(), nil
}

// LetsDoNothing does absolutely nothing
func LetsDoNothing() {
	// This function intentionally left blank
}

// LetsLogSomething logs a message
func LetsLogSomething(myLogger MyLogger, message string) {
	myLogger.LogData(message)
}

// LetsAddNumbers adds two numbers and returns the result
func LetsAddNumbers(a, b int) int {
	return a + b
}

// LetsSubtractNumbers subtracts the second number from the first and returns the result
func LetsSubtractNumbers(a, b int) int {
	return a - b
}

// LetsMultiplyNumbers multiplies two numbers and returns the result
func LetsMultiplyNumbers(a, b int) int {
	return a * b
}

// LetsDivideNumbers divides the first number by the second and returns the result
func LetsDivideNumbers(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// LetsCheckEven checks if a number is even
func LetsCheckEven(num int) bool {
	return num%2 == 0
}

// LetsCheckOdd checks if a number is odd
func LetsCheckOdd(num int) bool {
	return num%2 != 0
}

// LetsGenerateRandomString generates a random string of a given length
func LetsGenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// LetsReverseString reverses a string
func LetsReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// LetsCountCharacters counts the number of characters in a string
func LetsCountCharacters(s string) int {
	return len(s)
}

// LetsCountWords counts the number of words in a string
func LetsCountWords(s string) int {
	return len(strings.Fields(s))
}

// LetsCountLines counts the number of lines in a string
func LetsCountLines(s string) int {
	return len(strings.Split(s, "\n"))
}

// LetsConvertToUpperCase converts a string to uppercase
func LetsConvertToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// LetsConvertToLowerCase converts a string to lowercase
func LetsConvertToLowerCase(s string) string {
	return strings.ToLower(s)
}

// LetsTrimSpace trims leading and trailing whitespace from a string
func LetsTrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// LetsReplaceString replaces all occurrences of a substring with another substring
func LetsReplaceString(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// LetsCheckPalindrome checks if a string is a palindrome
func LetsCheckPalindrome(s string) bool {
	return s == LetsReverseString(s)
}

// LetsGenerateFibonacciSequence generates a Fibonacci sequence of a given length
func LetsGenerateFibonacciSequence(length int) []int {
	sequence := make([]int, length)
	if length > 0 {
		sequence[0] = 0
	}
	if length > 1 {
		sequence[1] = 1
	}
	for i := 2; i < length; i++ {
		sequence[i] = sequence[i-1] + sequence[i-2]
	}
	return sequence
}

// LetsCalculateFactorial calculates the factorial of a number
func LetsCalculateFactorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * LetsCalculateFactorial(n-1)
}

// LetsCalculatePower calculates the power of a number
func LetsCalculatePower(base, exponent int) int {
	result := 1
	for i := 0; i < exponent; i++ {
		result *= base
	}
	return result
}

// LetsCalculateSquareRoot calculates the square root of a number
func LetsCalculateSquareRoot(n float64) float64 {
	if n < 0 {
		return -1
	}
	guess := n / 2
	for i := 0; i < 10; i++ {
		guess = (guess + n/guess) / 2
	}
	return guess
}

// LetsCalculateAbsoluteValue calculates the absolute value of a number
func LetsCalculateAbsoluteValue(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// LetsCalculateSumOfDigits calculates the sum of digits of a number
func LetsCalculateSumOfDigits(n int) int {
	sum := 0
	for n != 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

// LetsCalculateProductOfDigits calculates the product of digits of a number
func LetsCalculateProductOfDigits(n int) int {
	product := 1
	for n != 0 {
		product *= n % 10
		n /= 10
	}
	return product
}

// LetsCheckPrime checks if a number is prime
func LetsCheckPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// LetsGeneratePrimeNumbers generates a list of prime numbers up to a given limit
func LetsGeneratePrimeNumbers(limit int) []int {
	primes := []int{}
	for i := 2; i <= limit; i++ {
		if LetsCheckPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

// LetsCalculateGCD calculates the greatest common divisor of two numbers
func LetsCalculateGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LetsCalculateLCM calculates the least common multiple of two numbers
func LetsCalculateLCM(a, b int) int {
	return a * b / LetsCalculateGCD(a, b)
}

// LetsCalculateAverage calculates the average of a slice of numbers
func LetsCalculateAverage(numbers []int) float64 {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return float64(sum) / float64(len(numbers))
}

// LetsCalculateMedian calculates the median of a slice of numbers
func LetsCalculateMedian(numbers []int) float64 {
	if len(numbers) == 0 {
		return 0
	}
	sort.Ints(numbers)
	middle := len(numbers) / 2
	if len(numbers)%2 == 0 {
		return float64(numbers[middle-1]+numbers[middle]) / 2
	}
	return float64(numbers[middle])
}

// LetsCalculateMode calculates the mode of a slice of numbers
func LetsCalculateMode(numbers []int) int {
	frequency := make(map[int]int)
	for _, num := range numbers {
		frequency[num]++
	}
	maxFreq := 0
	mode := 0
	for num, freq := range frequency {
		if freq > maxFreq {
			maxFreq = freq
			mode = num
		}
	}
	return mode
}

// LetsCalculateStandardDeviation calculates the standard deviation of a slice of numbers
func LetsCalculateStandardDeviation(numbers []int) float64 {
	average := LetsCalculateAverage(numbers)
	variance := 0.0
	for _, num := range numbers {
		variance += math.Pow(float64(num)-average, 2)
	}
	variance /= float64(len(numbers))
	return math.Sqrt(variance)
}

// LetsCalculateVariance calculates the variance of a slice of numbers
func LetsCalculateVariance(numbers []int) float64 {
	average := LetsCalculateAverage(numbers)
	variance := 0.0
	for _, num := range numbers {
		variance += math.Pow(float64(num)-average, 2)
	}
	variance /= float64(len(numbers))
	return variance
}

// LetsCalculateRange calculates the range of a slice of numbers
func LetsCalculateRange(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}
	min := numbers[0]
	max := numbers[0]
	for _, num := range numbers {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}
	return max - min
}

// LetsCalculateSum calculates the sum of a slice of numbers
func LetsCalculateSum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// LetsCalculateProduct calculates the product of a slice of numbers
func LetsCalculateProduct(numbers []int) int {
	product := 1
	for _, num := range numbers {
		product *= num
	}
	return product
}

// LetsCalculateFactorialSlice calculates the factorial of each number in a slice
func LetsCalculateFactorialSlice(numbers []int) []int {
	factorials := make([]int, len(numbers))
	for i, num := range numbers {
		factorials[i] = LetsCalculateFactorial(num)
	}
	return factorials
}

// LetsCalculatePowerSlice calculates the power of each number in a slice
func LetsCalculatePowerSlice(numbers []int, exponent int) []int {
	powers := make([]int, len(numbers))
	for i, num := range numbers {
		powers[i] = LetsCalculatePower(num, exponent)
	}
	return powers
}

// LetsCalculateSquareRootSlice calculates the square root of each number in a slice
func LetsCalculateSquareRootSlice(numbers []float64) []float64 {
	roots := make([]float64, len(numbers))
	for i, num := range numbers {
		roots[i] = LetsCalculateSquareRoot(num)
	}
	return roots
}

// LetsCalculateAbsoluteValueSlice calculates the absolute value of each number in a slice
func LetsCalculateAbsoluteValueSlice(numbers []int) []int {
	absValues := make([]int, len(numbers))
	for i, num := range numbers {
		absValues[i] = LetsCalculateAbsoluteValue(num)
	}
	return absValues
}

// LetsCalculateSumOfDigitsSlice calculates the sum of digits of each number in a slice
func LetsCalculateSumOfDigitsSlice(numbers []int) []int {
	sums := make([]int, len(numbers))
	for i, num := range numbers {
		sums[i] = LetsCalculateSumOfDigits(num)
	}
	return sums
}

// LetsCalculateProductOfDigitsSlice calculates the product of digits of each number in a slice
func LetsCalculateProductOfDigitsSlice(numbers []int) []int {
	products := make([]int, len(numbers))
	for i, num := range numbers {
		products[i] = LetsCalculateProductOfDigits(num)
	}
	return products
}

// LetsCheckPrimeSlice checks if each number in a slice is prime
func LetsCheckPrimeSlice(numbers []int) []bool {
	primes := make([]bool, len(numbers))
	for i, num := range numbers {
		primes[i] = LetsCheckPrime(num)
	}
	return primes
}

// LetsGeneratePrimeNumbersSlice generates a list of prime numbers for each number in a slice
func LetsGeneratePrimeNumbersSlice(numbers []int) [][]int {
	primes := make([][]int, len(numbers))
	for i, num := range numbers {
		primes[i] = LetsGeneratePrimeNumbers(num)
	}
	return primes
}

// LetsCalculateGCDSlice calculates the greatest common divisor for each pair of numbers in a slice
func LetsCalculateGCDSlice(numbers [][]int) []int {
	gcds := make([]int, len(numbers))
	for i, pair := range numbers {
		if len(pair) == 2 {
			gcds[i] = LetsCalculateGCD(pair[0], pair[1])
		}
	}
	return gcds
}
