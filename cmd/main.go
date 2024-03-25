package main

import (

  "github.com/LegendaryLlama37/api_consumer_concurrent/apiquery"
	//"encoding/json"
	"fmt"
	//"net/http"
	"os"
	"time"
)

const (
	instanceURL = "https://ServiceNowInstance.service-now.com/"
)

func main() {
	// Fetch credentials from environment variables
	username := os.Getenv("SN_USERNAME")
	password := os.Getenv("SN_PASSWORD")
	if username == "" || password == "" {
		fmt.Println("Please provide ServiceNow username and password as environment variables (SN_USERNAME and SN_PASSWORD)")
		return
	}

	// Fetch releases assigned to a user
	releases := getReleasesAssignedToUser("user1", username, password)
	for _, release := range releases.Result {
		fmt.Printf("Release Number: %s\n", release.Number)
		stories := getStoriesInRelease(release.SysID,  username, password)
		for _, story := range stories.Result {
			fmt.Printf(" Story Number: %s - %s\n", story.Number, story.Name)
		}
		fmt.Println()
	}

	// Get completed update sets within the last two weeks
	updateSets := getCompletedUpdateSetWithinTwoWeeks(username, password)
	fmt.Println("Completed Update Sets within the last two weeks:")
	for _, updateSet := range updateSets.Result {
		fmt.Printf("Update Set ID: %s\n", updateSet.SysID)
	}
}

func getReleasesAssignedToUser(user, username, password string) apiquery.Release {
	url := fmt.Sprintf("%s/api/now/table/rm_release_scrum?sysparm_query=assigned_to=%s", instanceURL, user)
	return queryServiceNowAPI(url, username, password)
}


func getStoriesInRelease(releaseID, username, password string) apiquery.Story {
	url := fmt.Sprintf("%s/api/now/table/rm_story?sysparm_query=releases=%s", instanceURL, releaseID)
	return queryServiceNowAPI(url, username, password)
}


func getCompletedUpdateSetWithinTwoWeeks(username, password string) apiquery.UpdateSet {
	now := time.Now()
	twoWeeksAgo := now.AddDate(0, 0, -14)
	formattedTwoWeeksAgo := twoWeeksAgo.Format("2006-01-02 15:04:05")

	url := fmt.Sprintf("%s/api/now/table/sys_update_set?sysparm_query=complete=true^sys_updated_on>=javascript:gs.dateGenerate('%s')", instanceURL, formattedTwoWeeksAgo)
	return queryServiceNowAPI(url, username, password)
}


func queryServiceNowAPI(url, username, password string) interface{} {
	credentials := apiquery.Credentials{APIKey: ""}
	results, err := apiquery.FetchData(url, credentials)
	if err != nil {
		fmt.Printf("Error fetching data from ServiceNow API: %v\n", err)
		return nil
	}
	return results
}

