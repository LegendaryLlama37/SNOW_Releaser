package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/LegendaryLlama37/api_consumer_concurrent/apiquery"
)

const (
	prodInstanceURL  = "https://prod-instance.service-now.com/api/now/table/release"
	devInstanceURL   = "https://dev-instance.service-now.com/api/now/table/update_set"
	apiKeyProd        = "your_production_api_key"
	apiKeyDev         = "your_dev_api_key"
	updateSetSearch   = ""
	updateSetAgeLimit = 14 * 24 * time.Hour         // Update set age limit (2 weeks)
)

func main() {
	// Define URLs and API keys for concurrent fetching
	urlAPIKeyMap := map[string]string{
		prodInstanceURL: apiKeyProd,
		devInstanceURL:  apiKeyDev,
	}

	// Fetch data from both instances concurrently
	dataMap, err := apiquery.FetchDataConcurrently(urlAPIKeyMap)
	if err != nil {
		log.Fatal(err)
	}

	// Extract releases from production data
	prodReleases, ok := dataMap[prodInstanceURL].([]map[string]interface{})
	if !ok {
		log.Fatal("Error parsing production release data")
	}

	// Loop through releases and search for update sets in dev instance
	for _, release := range prodReleases {
		releaseName := release["name"].(string)
		fmt.Printf("Searching update sets for release: %s\n", releaseName)

		// Extract story numbers from release data (assuming "short_description" contains story numbers)
		storyNumbers := extractStoryNumbers(release["short_description"].(string))

		updateSets := searchDevUpdateSets(storyNumbers, dataMap[devInstanceURL])
		if updateSets == nil {
			log.Printf("Error searching update sets for release %s\n", releaseName)
			continue
		}

		// Print update sets that match story numbers and are within age limit
		printMatchingUpdateSets(updateSets)
	}
}

func extractStoryNumbers(description string) []string {
	// Extract story numbers using a regular expression (STRY- format)
	storyNumberRegex := regexp.MustCompile(`STRY-[A-Z0-9]+`)
	matches := storyNumberRegex.FindAllString(description, -1)
	return matches
}

func searchDevUpdateSets(storyNumbers []string, devData interface{}) []map[string]interface{} {
	twoWeeksAgo := time.Now().Add(-updateSetAgeLimit)
	filter := fmt.Sprintf("sys_created_on>=%s^opened^released", twoWeeksAgo.Format("2006-01-02"))

	// Iterate through dev update sets and filter based on story numbers and criteria
	var matchingSets []map[string]interface{}
	for _, updateSet := range devData.([]map[string]interface{}) {
		updateSetName := updateSet["name"].(string)
		if updateSet["sys_created_on"].(string) >= twoWeeksAgo.Format("2006-01-02T15:04:05") {
			for _, storyNumber := range storyNumbers {
				if strings.Contains(updateSetName, storyNumber) {
					matchingSets = append(matchingSets, updateSet)
					break // Move to next update set after finding a match
				}
			}
		}
	}
	return matchingSets
}

func printMatchingUpdateSets(updateSets []map[string]interface{}) {
	for _, updateSet := range updateSets {
		updateSetName := updateSet["name"].(string)
		fmt.Printf("  - Update Set: %s\n", updateSetName)
	}
}

