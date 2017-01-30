# jiraintegration

This tool is designed to easy test cucumber scenarios definew within Jira Xray Integration.

##Jira Integration
Download/upload test/results to Jira server using Xray plugin.

* **jiraintegration -download_test** *Download tests from Jira server*
  * **Mandatory options**
    * -host 
    * -user
    * [-jiraIssues | -keys | -filter]
  * **Optional options**
    * -featuresDir
    * [-password | -passworFile]

* **jiraintegration -export** *Upload test results to Jira server (Json format mandatory)*
  * **Mandatory options**
    * -host 
    * -user
    * -resultFile
  * **Optional options**
    * [-password | -passworFile]

##Execute cucumber tests
Execute & get undefined steps usign ruby cucumber tool.

* **jiraintegration [-execute]** *Execute cucumber tests (Needs ruby dependencies)*
  * **Optional options**
    * -featuresDir
    * -resultFile
    * -format
    * -stepsRepo
    
* **jiraintegration -getUndefinedSteps** *Get undefined steps and place them under features/step_definitions/pending_xxx.rb*
  * **Optional options**a
    * -featuresDir
    * -stepsDir	
    * -stepsRepo
  
##All in one command 
Download, get undefined steps, execute & upload results to Jira server within 1 command

* **jiraintegration -executeRemote** *Download tests from Jira server*
  * **Mandatory options**
    * -host 
    * -user
    * [-jiraIssues | -keys | -filter]
  * **Optional options**
    * -featuresDir
    * -resultFile
    * -stepsRepo
    * [-password | -passworFile]


##Command HELP

**jiraintegration -h**

```Bash
 -download_test
        Download test Scenarios from Jira
  -execute
        Execute test Scenarios from $featuresDir
  -executeRemote
        Download, Execute & Upload test from/to $host server
  -export
        Export cucumber results to Jira
  -featuresDir string
        Target directory for downloaded tests (default "./features")
  -filter int
        Filter query to retrieve tests
  -format string
        Cucumber test result format
  -getUndefinedSteps
        Get undefined steps definition for provided features
  -host string
        Jira server URL
  -jiraIssues string
         Jira-keys list of issues with associated test separated by ';'
  -keys string
         Jira-keys list of test-set Issues separated by ';'
  -password string
        Jira server password
  -passwordFile string
        Jira server password file
  -resultFile string
         Path to cucumber result file
  -stepsDir string
        Target directory for place unddefined steps generated file (default "./features/step_definitions")
  -stepsRepo string
        Url to repository containing steps definitions (cucumber structure needed in the repo root=/features)
  -user string
        Jira server use
```

