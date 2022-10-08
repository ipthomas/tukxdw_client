// tukpdq provides a golang implementtion of, IHE PIXm,IHE PIXv3 and IHE PDQv3 Client Consumers
//
// There is currently no authentication implemented. The func (i *PDQQuery) newRequest() error is used to handle the http request/response and should be amended according to your authentication requirements
//
// Struct PDQQuery implements the tukpdq.PDQ() interface
//
//	type PDQQuery struct {
//		Count        int          `json:"count"`
//		PID          string       `json:"pid"`
//		PIDOID       string       `json:"pidoid"`
//		PIX_URL      string       `json:"pixurl"`
//		NHS_OID      string       `json:"nhsoid"`
//		Region_OID   string       `json:"regionoid"`
//		Timeout      int64        `json:"timeout"`
//		StatusCode   int          `json:"statuscode"`
//		Response     []byte       `json:"response"`
//		PIXmResponse PIXmResponse `json:"pixmresponse"`
//		Patients     []PIXPatient `json:"patients"`
//	}
//
//	 PID is the MRN or NHS ID or Regional/xds ID and is required
//	 Region_OID is the Regional/XDS OID and is required
//	 PIX_URL is the PIXm WS end point and is required.
//	 PID_OID is required if the PID is not an NHS ID. If pid length = 10 and no PID_OID is provided, the pid is assumed to be a NHS ID and the PID_OID is set to the NHS offical NHS ID OID (2.16.840.1.113883.2.1.4.1)
//	 Timeout sets the http context timeout in seconds and is optional. Default is 5 secs
//	 NHS_OID is optional. Default is 2.16.840.1.113883.2.1.4.1
//	 Count will be set from the pixm response to the number of patients found matching the query
//	 Response will contain the PIXm response in []byte format
//	 PIXmResponse will contain the initialised PIXmResponse struc from the Response []byte
//	 StatusCode will be set from the PIXm Server http response header statuscode
//	 []Patients is any array of PIXPatient structs containing all matched patients. Hopefully just 1 !!
//
//	Example usage:
//		pdq := tukpdq.PIXmQuery{
//			PID:        "9999999468",
//			Region_OID: "2.16.840.1.113883.2.1.3.31.2.1.1",
//			PIX_URL:    "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient",
//		}
//		if err = tukpdq.PDQ(&pdq); err == nil {
//			log.Printf("Patient %s %s is registered", pdq.Patients[0].GivenName, pdq.Patients[0].FamilyName)
//		} else {
//			log.Println(err.Error())
//		}
//
//	Running the above example produces the following Log output:
//
//	2022/09/12 14:02:55.510679 tukpdq.go:188: HTTP GET Request Headers
//
//	2022/09/12 14:02:55.510834 tukpdq.go:190: {
//	  "Accept": [
//	    "*/*"
//	  ],
//	  "Connection": [
//	    "keep-alive"
//	  ],
//	  "Content-Type": [
//	    "application/json"
//	  ]
//	}
//
// 2022/09/12 14:02:55.510860 tukpdq.go:191: HTTP Request
// URL = http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient?identifier=2.16.840.1.113883.2.1.4.1%7C9999999468&_format=json&_pretty=true
// 2022/09/12 14:02:55.851605 tukpdq.go:194: HTML Response - Status Code = 200
//
//	{
//	  "resourceType": "Bundle",
//	  "id": "53c44d32-fb2c-4dfb-b819-db2150e6fa87",
//	  "type": "searchset",
//	  "total": 1,
//	  "link": [ {
//	    "relation": "self",
//	    "url": "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient?_format=json&_pretty=true&identifier=2.16.840.1.113883.2.1.4.1%7C9999999468"
//	  } ],
//	  "entry": [ {
//	    "fullUrl": "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient/VFNVSy4xNjYxOTc2MjMwMjYxMSYyLjE2Ljg0MC4xLjExMzg4My4yLjEuMy4zMS4yLjEuMS4xLjMuMS4x",
//	    "resource": {
//	      "resourceType": "Patient",
//	      "id": "VFNVSy4xNjYxOTc2MjMwMjYxMSYyLjE2Ljg0MC4xLjExMzg4My4yLjEuMy4zMS4yLjEuMS4xLjMuMS4x",
//	      "extension": [ {
//	        "url": "http://hl7.org/fhir/StructureDefinition/patient-citizenship",
//	        "valueCodeableConcept": {
//	          "coding": [ {
//	            "code": "GBR"
//	          } ]
//	        }
//	      }, {
//	        "url": "http://hl7.org/fhir/StructureDefinition/patient-nationality",
//	        "valueCodeableConcept": {
//	          "coding": [ {
//	            "code": "GBR"
//	          } ]
//	        }
//	      } ],
//	      "identifier": [ {
//	        "system": "urn:oid:2.16.840.1.113883.2.1.4.1",
//	        "value": "9999999468"
//	      }, {
//	        "use": "usual",
//	        "system": "urn:oid:2.16.840.1.113883.2.1.3.31.2.1.1.1.3.1.1",
//	        "value": "TSUK.16619762302611"
//	      }, {
//	        "system": "urn:oid:2.16.840.1.113883.2.1.3.31.2.1.1",
//	        "value": "REG.1MWU5C92M2"
//	      } ],
//	      "active": true,
//	      "name": [ {
//	        "use": "official",
//	        "family": "Testpatient",
//	        "given": [ "Nhs" ]
//	      } ],
//	      "telecom": [ {
//	        "system": "phone",
//	        "value": "07777661324",
//	        "use": "work"
//	      }, {
//	        "system": "email",
//	        "value": "nhs.testpatient@nhs.net",
//	        "use": "work"
//	      } ],
//	      "gender": "male",
//	      "birthDate": "1962-04-04",
//	      "address": [ {
//	        "line": [ "Preston Road" ],
//	        "city": "Preston",
//	        "state": "Lancashire",
//	        "postalCode": "PR1 1PR",
//	        "country": "GBR"
//	      } ],
//	      "maritalStatus": {
//	        "coding": [ {
//	          "code": "D"
//	        } ]
//	      },
//	      "multipleBirthBoolean": false
//	    }
//	  } ]
//	}
//
// 2022/09/12 14:02:55.852334 tukpdq.go:102: 1 Patient Entries in Response
// 2022/09/12 14:02:55.852392 tukpdq.go:122: Set NHS ID 9999999468 2.16.840.1.113883.2.1.4.1
// 2022/09/12 14:02:55.852427 tukpdq.go:117: Set PID TSUK.16619762302611 2.16.840.1.113883.2.1.3.31.2.1.1.1.3.1.1
// 2022/09/12 14:02:55.852455 tukpdq.go:112: Set Reg ID REG.1MWU5C92M2 2.16.840.1.113883.2.1.3.31.2.1.1
// 2022/09/12 14:02:55.852546 tukpdq.go:149: Added Patient 9999999468 to response
// 2022/09/12 14:02:55.852569 main.go:84: Patient Nhs Testpatient is registered
package tukpdq

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukhttp"
	"github.com/ipthomas/tukutil"
)

type PDQQuery struct {
	Server_Mode     string           `json:",omitempty"`
	Server_URL      string           `json:",omitempty"`
	CGL_X_Api_Key   string           `json:",omitempty"`
	NHS_ID          string           `json:",omitempty"`
	NHS_OID         string           `json:",omitempty"`
	MRN_ID          string           `json:",omitempty"`
	MRN_OID         string           `json:",omitempty"`
	REG_ID          string           `json:",omitempty"`
	REG_OID         string           `json:",omitempty"`
	Timeout         int64            `json:",omitempty"`
	Cache           bool             `json:",omitempty"`
	Used_PID        string           `json:",omitempty"`
	Used_PID_OID    string           `json:",omitempty"`
	Request         []byte           `json:",omitempty"`
	Response        []byte           `json:",omitempty"`
	StatusCode      int              `json:",omitempty"`
	Count           int              `json:",omitempty"`
	PDQv3Response   *PDQv3Response   `json:",omitempty"`
	PIXv3Response   *PIXv3Response   `json:",omitempty"`
	PIXmResponse    *PIXmResponse    `json:",omitempty"`
	Patients        *[]TUKPatient    `json:",omitempty"`
	CGLUserResponse *CGLUserResponse `json:",omitempty"`
}
type CGLUserResponse struct {
	Data struct {
		Client struct {
			BasicDetails struct {
				Address struct {
					AddressLine1 string `json:"addressLine1,omitempty"`
					AddressLine2 string `json:"addressLine2,omitempty"`
					AddressLine3 string `json:"addressLine3,omitempty"`
					AddressLine4 string `json:"addressLine4,omitempty"`
					AddressLine5 string `json:"addressLine5,omitempty"`
					PostCode     string `json:"postCode,omitempty"`
				} `json:"address,omitempty"`
				BirthDate                    string `json:"birthDate,omitempty"`
				Disability                   string `json:"disability,omitempty"`
				LastEngagementByCGLDate      string `json:"lastEngagementByCGLDate,omitempty"`
				LastFaceToFaceEngagementDate string `json:"lastFaceToFaceEngagementDate,omitempty"`
				LocalIdentifier              int    `json:"localIdentifier,omitempty"`
				Name                         struct {
					Family string `json:"family,omitempty"`
					Given  string `json:"given,omitempty"`
				} `json:"name,omitempty"`
				NextCGLAppointmentDate string `json:"nextCGLAppointmentDate,omitempty"`
				NhsNumber              string `json:"nhsNumber,omitempty"`
				SexAtBirth             string `json:"sexAtBirth,omitempty"`
			} `json:"basicDetails,omitempty"`
			BbvInformation struct {
				BbvTested        string `json:"bbvTested,omitempty"`
				HepCLastTestDate string `json:"hepCLastTestDate,omitempty"`
				HepCResult       string `json:"hepCResult,omitempty"`
				HivPositive      string `json:"hivPositive,omitempty"`
			} `json:"bbvInformation,omitempty"`
			DrugTestResults struct {
				DrugTestDate          string `json:"drugTestDate,omitempty"`
				DrugTestSample        string `json:"drugTestSample,omitempty"`
				DrugTestStatus        string `json:"drugTestStatus,omitempty"`
				InstantOrConfirmation string `json:"instantOrConfirmation,omitempty"`
				Results               struct {
					Amphetamine     string `json:"amphetamine,omitempty"`
					Benzodiazepine  string `json:"benzodiazepine,omitempty"`
					Buprenorphine   string `json:"buprenorphine,omitempty"`
					Cannabis        string `json:"cannabis,omitempty"`
					Cocaine         string `json:"cocaine,omitempty"`
					Eddp            string `json:"eddp,omitempty"`
					Fentanyl        string `json:"fentanyl,omitempty"`
					Ketamine        string `json:"ketamine,omitempty"`
					Methadone       string `json:"methadone,omitempty"`
					Methamphetamine string `json:"methamphetamine,omitempty"`
					Morphine        string `json:"morphine,omitempty"`
					Opiates         string `json:"opiates,omitempty"`
					SixMam          string `json:"sixMam,omitempty"`
					Tramadol        string `json:"tramadol,omitempty"`
				} `json:"results,omitempty"`
			} `json:"drugTestResults,omitempty"`
			PrescribingInformation []string `json:"prescribingInformation,omitempty"`
			RiskInformation        struct {
				LastSelfReportedDate string `json:"lastSelfReportedDate,omitempty"`
				MentalHealthDomain   struct {
					AttemptedSuicide                            string `json:"attemptedSuicide,omitempty"`
					CurrentOrPreviousSelfHarm                   string `json:"currentOrPreviousSelfHarm,omitempty"`
					DiagnosedMentalHealthCondition              string `json:"diagnosedMentalHealthCondition,omitempty"`
					FrequentLifeThreateningSelfHarm             string `json:"frequentLifeThreateningSelfHarm,omitempty"`
					Hallucinations                              string `json:"hallucinations,omitempty"`
					HospitalAdmissionsForMentalHealth           string `json:"hospitalAdmissionsForMentalHealth,omitempty"`
					NoIdentifiedRisk                            string `json:"noIdentifiedRisk,omitempty"`
					NotEngagingWithSupport                      string `json:"notEngagingWithSupport,omitempty"`
					NotTakingPrescribedMedicationAsInstructed   string `json:"notTakingPrescribedMedicationAsInstructed,omitempty"`
					PsychiatricOrPreviousCrisisTeamIntervention string `json:"psychiatricOrPreviousCrisisTeamIntervention,omitempty"`
					Psychosis                                   string `json:"psychosis,omitempty"`
					SelfReportedMentalHealthConcerns            string `json:"selfReportedMentalHealthConcerns,omitempty"`
					ThoughtsOfSuicideOrSelfHarm                 string `json:"thoughtsOfSuicideOrSelfHarm,omitempty"`
				} `json:"mentalHealthDomain,omitempty"`
				RiskOfHarmToSelfDomain struct {
					AssessedAsNotHavingMentalCapacity  string `json:"assessedAsNotHavingMentalCapacity,omitempty"`
					BeliefTheyAreWorthless             string `json:"beliefTheyAreWorthless,omitempty"`
					Hoarding                           string `json:"hoarding,omitempty"`
					LearningDisability                 string `json:"learningDisability,omitempty"`
					MeetsSafeguardingAdultsThreshold   string `json:"meetsSafeguardingAdultsThreshold,omitempty"`
					NoIdentifiedRisk                   string `json:"noIdentifiedRisk,omitempty"`
					OngoingConcernsRelatingToOwnSafety string `json:"ongoingConcernsRelatingToOwnSafety,omitempty"`
					ProblemsMaintainingPersonalHygiene string `json:"problemsMaintainingPersonalHygiene,omitempty"`
					ProblemsMeetingNutritionalNeeds    string `json:"problemsMeetingNutritionalNeeds,omitempty"`
					RequiresIndependentAdvocacy        string `json:"requiresIndependentAdvocacy,omitempty"`
					SelfNeglect                        string `json:"selfNeglect,omitempty"`
				} `json:"riskOfHarmToSelfDomain,omitempty"`
				SocialDomain struct {
					FinancialProblems         string `json:"financialProblems,omitempty"`
					HomelessRoughSleepingNFA  string `json:"homelessRoughSleepingNFA,omitempty"`
					HousingAtRisk             string `json:"housingAtRisk,omitempty"`
					NoIdentifiedRisk          string `json:"noIdentifiedRisk,omitempty"`
					SociallyIsolatedNoSupport string `json:"sociallyIsolatedNoSupport,omitempty"`
				} `json:"socialDomain,omitempty"`
				SubstanceMisuseDomain struct {
					ConfusionOrDisorientation string `json:"ConfusionOrDisorientation,omitempty"`
					AdmissionToAandE          string `json:"admissionToAandE,omitempty"`
					BlackoutOrSeizures        string `json:"blackoutOrSeizures,omitempty"`
					ConcurrentUse             string `json:"concurrentUse,omitempty"`
					HigherRiskDrinking        string `json:"higherRiskDrinking,omitempty"`
					InjectedByOthers          string `json:"injectedByOthers,omitempty"`
					Injecting                 string `json:"injecting,omitempty"`
					InjectingInNeckOrGroin    string `json:"injectingInNeckOrGroin,omitempty"`
					NoIdentifiedRisk          string `json:"noIdentifiedRisk,omitempty"`
					PolyDrugUse               string `json:"polyDrugUse,omitempty"`
					PreviousOverDose          string `json:"previousOverDose,omitempty"`
					RecentPrisonRelease       string `json:"recentPrisonRelease,omitempty"`
					ReducedTolerance          string `json:"reducedTolerance,omitempty"`
					SharingWorks              string `json:"sharingWorks,omitempty"`
					Speedballing              string `json:"speedballing,omitempty"`
					UsingOnTop                string `json:"usingOnTop,omitempty"`
				} `json:"substanceMisuseDomain,omitempty"`
			} `json:"riskInformation,omitempty"`
			SafeguardingInformation struct {
				LastReviewDate     string `json:"lastReviewDate,omitempty"`
				RiskHarmFromOthers string `json:"riskHarmFromOthers,omitempty"`
				RiskToAdults       string `json:"riskToAdults,omitempty"`
				RiskToChildrenOrYP string `json:"riskToChildrenOrYP,omitempty"`
				RiskToSelf         string `json:"riskToSelf,omitempty"`
			} `json:"safeguardingInformation,omitempty"`
		} `json:"client,omitempty"`
		KeyWorker struct {
			LocalIdentifier int `json:"localIdentifier,omitempty"`
			Name            struct {
				Family string `json:"family,omitempty"`
				Given  string `json:"given,omitempty"`
			} `json:"name"`
			Telecom string `json:"telecom,omitempty"`
		} `json:"keyWorker,omitempty"`
	} `json:"data,omitempty"`
}
type PDQv3Response struct {
	XMLName xml.Name `xml:"Envelope"`
	S       string   `xml:"S,attr"`
	Env     string   `xml:"env,attr"`
	Header  struct {
		Action struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"Action"`
		MessageID struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"MessageID"`
		RelatesTo struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"RelatesTo"`
		To struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"To"`
	} `xml:"Header"`
	Body struct {
		PRPAIN201306UV02 struct {
			Xmlns      string `xml:"xmlns,attr"`
			ITSVersion string `xml:"ITSVersion,attr"`
			ID         struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"id"`
			CreationTime struct {
				Value string `xml:"value,attr"`
			} `xml:"creationTime"`
			VersionCode struct {
				Code string `xml:"code,attr"`
			} `xml:"versionCode"`
			InteractionId struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"interactionId"`
			ProcessingCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingCode"`
			ProcessingModeCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingModeCode"`
			AcceptAckCode struct {
				Code string `xml:"code,attr"`
			} `xml:"acceptAckCode"`
			Receiver struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
						Root                   string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								Root                   string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"receiver"`
			Sender struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						Root string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								Root string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"sender"`
			Acknowledgement struct {
				TypeCode struct {
					Code string `xml:"code,attr"`
				} `xml:"typeCode"`
				TargetMessage struct {
					ID struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"id"`
				} `xml:"targetMessage"`
			} `xml:"acknowledgement"`
			ControlActProcess struct {
				ClassCode string `xml:"classCode,attr"`
				MoodCode  string `xml:"moodCode,attr"`
				Code      struct {
					Code       string `xml:"code,attr"`
					CodeSystem string `xml:"codeSystem,attr"`
				} `xml:"code"`
				Subject struct {
					ContextConductionInd string `xml:"contextConductionInd,attr"`
					TypeCode             string `xml:"typeCode,attr"`
					RegistrationEvent    struct {
						ClassCode string `xml:"classCode,attr"`
						MoodCode  string `xml:"moodCode,attr"`
						ID        struct {
							NullFlavor string `xml:"nullFlavor,attr"`
						} `xml:"id"`
						StatusCode struct {
							Code string `xml:"code,attr"`
						} `xml:"statusCode"`
						Subject1 struct {
							TypeCode string `xml:"typeCode,attr"`
							Patient  struct {
								ClassCode string `xml:"classCode,attr"`
								ID        []struct {
									AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
									Extension              string `xml:"extension,attr"`
									Root                   string `xml:"root,attr"`
								} `xml:"id"`
								StatusCode struct {
									Code string `xml:"code,attr"`
								} `xml:"statusCode"`
								EffectiveTime struct {
									Value string `xml:"value,attr"`
								} `xml:"effectiveTime"`
								PatientPerson struct {
									ClassCode      string `xml:"classCode,attr"`
									DeterminerCode string `xml:"determinerCode,attr"`
									Name           struct {
										Use    string `xml:"use,attr"`
										Given  string `xml:"given"`
										Family string `xml:"family"`
									} `xml:"name"`
									AdministrativeGenderCode struct {
										Code           string `xml:"code,attr"`
										CodeSystem     string `xml:"codeSystem,attr"`
										CodeSystemName string `xml:"codeSystemName,attr"`
									} `xml:"administrativeGenderCode"`
									BirthTime struct {
										Value string `xml:"value,attr"`
									} `xml:"birthTime"`
									DeceasedInd struct {
										Value string `xml:"value,attr"`
									} `xml:"deceasedInd"`
									MultipleBirthInd struct {
										Value string `xml:"value,attr"`
									} `xml:"multipleBirthInd"`
									Addr struct {
										StreetAddressLine string `xml:"streetAddressLine"`
										City              string `xml:"city"`
										State             string `xml:"state"`
										PostalCode        string `xml:"postalCode"`
										Country           string `xml:"country"`
									} `xml:"addr"`
									BirthPlace struct {
										Addr struct {
											City string `xml:"city"`
										} `xml:"addr"`
									} `xml:"birthPlace"`
								} `xml:"patientPerson"`
								ProviderOrganization struct {
									ClassCode      string `xml:"classCode,attr"`
									DeterminerCode string `xml:"determinerCode,attr"`
									ID             struct {
										NullFlavor string `xml:"nullFlavor,attr"`
									} `xml:"id"`
									ContactParty struct {
										ClassCode string `xml:"classCode,attr"`
									} `xml:"contactParty"`
								} `xml:"providerOrganization"`
								SubjectOf1 struct {
									TypeCode              string `xml:"typeCode,attr"`
									QueryMatchObservation struct {
										ClassCode string `xml:"classCode,attr"`
										MoodCode  string `xml:"moodCode,attr"`
										Code      struct {
											Code       string `xml:"code,attr"`
											CodeSystem string `xml:"codeSystem,attr"`
										} `xml:"code"`
										Value struct {
											Value string `xml:"value,attr"`
											Xsi   string `xml:"xsi,attr"`
											Type  string `xml:"type,attr"`
										} `xml:"value"`
									} `xml:"queryMatchObservation"`
								} `xml:"subjectOf1"`
							} `xml:"patient"`
						} `xml:"subject1"`
						Custodian struct {
							TypeCode       string `xml:"typeCode,attr"`
							AssignedEntity struct {
								ClassCode string `xml:"classCode,attr"`
								ID        struct {
									NullFlavor string `xml:"nullFlavor,attr"`
								} `xml:"id"`
							} `xml:"assignedEntity"`
						} `xml:"custodian"`
					} `xml:"registrationEvent"`
				} `xml:"subject"`
				QueryAck struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					QueryResponseCode struct {
						Code string `xml:"code,attr"`
					} `xml:"queryResponseCode"`
					ResultTotalQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultTotalQuantity"`
					ResultCurrentQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultCurrentQuantity"`
					ResultRemainingQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultRemainingQuantity"`
				} `xml:"queryAck"`
				QueryByParameter struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					ResponseModalityCode struct {
						Code string `xml:"code,attr"`
					} `xml:"responseModalityCode"`
					ResponsePriorityCode struct {
						Code string `xml:"code,attr"`
					} `xml:"responsePriorityCode"`
					MatchCriterionList string `xml:"matchCriterionList"`
					ParameterList      struct {
						LivingSubjectId struct {
							Value struct {
								Extension string `xml:"extension,attr"`
							} `xml:"value"`
							SemanticsText string `xml:"semanticsText"`
						} `xml:"livingSubjectId"`
					} `xml:"parameterList"`
				} `xml:"queryByParameter"`
			} `xml:"controlActProcess"`
		} `xml:"PRPA_IN201306UV02"`
	} `xml:"Body"`
}
type PIXv3Response struct {
	XMLName xml.Name `xml:"Envelope"`
	S       string   `xml:"S,attr"`
	Env     string   `xml:"env,attr"`
	Header  struct {
		Action struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"Action"`
		MessageID struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"MessageID"`
		RelatesTo struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"RelatesTo"`
		To struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"To"`
	} `xml:"Header"`
	Body struct {
		PRPAIN201310UV02 struct {
			Xmlns      string `xml:"xmlns,attr"`
			ITSVersion string `xml:"ITSVersion,attr"`
			ID         struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"id"`
			CreationTime struct {
				Value string `xml:"value,attr"`
			} `xml:"creationTime"`
			VersionCode struct {
				Code string `xml:"code,attr"`
			} `xml:"versionCode"`
			InteractionId struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"interactionId"`
			ProcessingCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingCode"`
			ProcessingModeCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingModeCode"`
			AcceptAckCode struct {
				Code string `xml:"code,attr"`
			} `xml:"acceptAckCode"`
			Receiver struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
						Root                   string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								Root                   string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"receiver"`
			Sender struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						Root string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								Root string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"sender"`
			Acknowledgement struct {
				TypeCode struct {
					Code string `xml:"code,attr"`
				} `xml:"typeCode"`
				TargetMessage struct {
					ID struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"id"`
				} `xml:"targetMessage"`
			} `xml:"acknowledgement"`
			ControlActProcess struct {
				ClassCode string `xml:"classCode,attr"`
				MoodCode  string `xml:"moodCode,attr"`
				Code      struct {
					Code       string `xml:"code,attr"`
					CodeSystem string `xml:"codeSystem,attr"`
				} `xml:"code"`
				Subject struct {
					TypeCode          string `xml:"typeCode,attr"`
					RegistrationEvent struct {
						ClassCode string `xml:"classCode,attr"`
						MoodCode  string `xml:"moodCode,attr"`
						ID        struct {
							NullFlavor string `xml:"nullFlavor,attr"`
						} `xml:"id"`
						StatusCode struct {
							Code string `xml:"code,attr"`
						} `xml:"statusCode"`
						Subject1 struct {
							TypeCode string `xml:"typeCode,attr"`
							Patient  struct {
								ClassCode string `xml:"classCode,attr"`
								ID        []struct {
									Extension              string `xml:"extension,attr"`
									Root                   string `xml:"root,attr"`
									AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								} `xml:"id"`
								StatusCode struct {
									Code string `xml:"code,attr"`
								} `xml:"statusCode"`
								PatientPerson struct {
									ClassCode      string `xml:"classCode,attr"`
									DeterminerCode string `xml:"determinerCode,attr"`
									Name           struct {
										Given  string `xml:"given"`
										Family string `xml:"family"`
									} `xml:"name"`
								} `xml:"patientPerson"`
							} `xml:"patient"`
						} `xml:"subject1"`
						Custodian struct {
							TypeCode       string `xml:"typeCode,attr"`
							AssignedEntity struct {
								ClassCode string `xml:"classCode,attr"`
								ID        struct {
									Root string `xml:"root,attr"`
								} `xml:"id"`
							} `xml:"assignedEntity"`
						} `xml:"custodian"`
					} `xml:"registrationEvent"`
				} `xml:"subject"`
				QueryAck struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					QueryResponseCode struct {
						Code string `xml:"code,attr"`
					} `xml:"queryResponseCode"`
					ResultTotalQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultTotalQuantity"`
					ResultCurrentQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultCurrentQuantity"`
					ResultRemainingQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultRemainingQuantity"`
				} `xml:"queryAck"`
				QueryByParameter struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					ResponsePriorityCode struct {
						Code string `xml:"code,attr"`
					} `xml:"responsePriorityCode"`
					ParameterList struct {
						PatientIdentifier struct {
							Value struct {
								AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								Extension              string `xml:"extension,attr"`
								Root                   string `xml:"root,attr"`
							} `xml:"value"`
							SemanticsText string `xml:"semanticsText"`
						} `xml:"patientIdentifier"`
					} `xml:"parameterList"`
				} `xml:"queryByParameter"`
			} `xml:"controlActProcess"`
		} `xml:"PRPA_IN201310UV02"`
	} `xml:"Body"`
}
type PIXmResponse struct {
	ResourceType string `json:"resourceType"`
	ID           string `json:"id"`
	Type         string `json:"type"`
	Total        int    `json:"total"`
	Link         []struct {
		Relation string `json:"relation"`
		URL      string `json:"url"`
	} `json:"link"`
	Entry []struct {
		FullURL  string `json:"fullUrl"`
		Resource struct {
			ResourceType string `json:"resourceType"`
			ID           string `json:"id"`
			Identifier   []struct {
				Use    string `json:"use,omitempty"`
				System string `json:"system"`
				Value  string `json:"value"`
			} `json:"identifier"`
			Active bool `json:"active"`
			Name   []struct {
				Use    string   `json:"use"`
				Family string   `json:"family"`
				Given  []string `json:"given"`
			} `json:"name"`
			Gender    string `json:"gender"`
			BirthDate string `json:"birthDate"`
			Address   []struct {
				Use        string   `json:"use"`
				Line       []string `json:"line"`
				City       string   `json:"city"`
				PostalCode string   `json:"postalCode"`
				Country    string   `json:"country"`
			} `json:"address"`
		} `json:"resource"`
	} `json:"entry"`
}
type TUKPatient struct {
	PIDOID     string `json:"pidoid"`
	PID        string `json:"pid"`
	REGOID     string `json:"regoid"`
	REGID      string `json:"regid"`
	NHSOID     string `json:"nhsoid"`
	NHSID      string `json:"nhsid"`
	GivenName  string `json:"givenname"`
	FamilyName string `json:"familyname"`
	Gender     string `json:"gender"`
	BirthDate  string `json:"birthdate"`
	Street     string `json:"street"`
	Town       string `json:"town"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Zip        string `json:"zip"`
}
type PDQInterface interface {
	pdq() error
}

var (
	pat_cache = make(map[string][]byte)
	pdq_cache = make(map[string]*[]TUKPatient)
)

func New_Transaction(i PDQInterface) error {
	return i.pdq()
}
func (i *PDQQuery) pdq() error {
	if err := i.setPDQ_ID(); err != nil {
		return err
	}
	return i.setPatient()
}
func (i *PDQQuery) setPDQ_ID() error {
	if i.Server_URL == "" {
		return errors.New("invalid request - pix server url is not set")
	}
	if i.REG_OID == "" {
		return errors.New("invalid request - reg oid is not set")
	}
	if i.Timeout == 0 {
		i.Timeout = 5
	}
	if i.NHS_OID == "" {
		i.NHS_OID = tukcnst.NHS_OID_DEFAULT
	}
	if i.MRN_ID != "" && i.MRN_OID != "" {
		i.Used_PID = i.MRN_ID
		i.Used_PID_OID = i.MRN_OID
	} else {
		if i.NHS_ID != "" {
			i.Used_PID = i.NHS_ID
			i.Used_PID_OID = i.NHS_OID
		} else {
			if i.REG_ID != "" && i.REG_OID != "" {
				i.Used_PID = i.REG_ID
				i.Used_PID_OID = i.REG_OID
			}
		}
	}
	if i.Used_PID == "" || i.Used_PID_OID == "" {
		return errors.New("invalid request - no suitable id and oid input values found which can be used for pdq query")
	}
	return nil
}
func (i *PDQQuery) setPatient() error {
	if i.Cache && i.Server_Mode != tukcnst.PDQ_SERVER_TYPE_CGL {
		if _, ok := pat_cache[i.Used_PID]; ok {
			log.Printf("Cache entry found for Patient ID %s", i.Used_PID)
			i.StatusCode = http.StatusOK
			i.Response = pat_cache[i.Used_PID]
			i.Patients = pdq_cache[i.Used_PID]
			i.Count = len(*i.Patients)
			return nil
		}
	}
	var tmplt *template.Template
	var err error
	i.StatusCode = http.StatusOK
	switch i.Server_Mode {
	case tukcnst.PDQ_SERVER_TYPE_CGL:
		i.Request = []byte(i.Server_URL + i.NHS_ID)
		httpReq := tukhttp.CGLRequest{
			Request:   i.Server_URL + i.NHS_ID,
			X_Api_Key: i.CGL_X_Api_Key,
		}
		if err = tukhttp.NewRequest(&httpReq); err == nil {
			if httpReq.StatusCode == http.StatusOK {
				json.Unmarshal(httpReq.Response, &i.CGLUserResponse)
				i.Count = 1
			}
		}
		i.Response = httpReq.Response
		i.StatusCode = httpReq.StatusCode
	case tukcnst.PDQ_SERVER_TYPE_IHE_PIXV3:
		if tmplt, err = template.New(tukcnst.PDQ_SERVER_TYPE_IHE_PIXV3).Funcs(tukutil.TemplateFuncMap()).Parse(tukcnst.GO_Template_PIX_V3_Request); err == nil {
			var b bytes.Buffer
			if err = tmplt.Execute(&b, i); err == nil {
				i.Request = b.Bytes()
				if err = i.newIHESOAPRequest(tukcnst.SOAP_ACTION_PIXV3_Request); err == nil {
					if err = xml.Unmarshal(i.Response, &i.PIXv3Response); err == nil {
						if i.PIXv3Response.Body.PRPAIN201310UV02.Acknowledgement.TypeCode.Code != "AA" {
							err = errors.New("acknowledgement code not equal aa, received " + i.PIXv3Response.Body.PRPAIN201310UV02.Acknowledgement.TypeCode.Code)
						} else {
							i.Count, _ = strconv.Atoi(i.PIXv3Response.Body.PRPAIN201310UV02.ControlActProcess.QueryAck.ResultTotalQuantity.Value)
							if i.Count > 0 {
								pat := TUKPatient{
									PIDOID: i.MRN_OID,
									PID:    i.MRN_ID,
									REGOID: i.REG_OID,
									REGID:  i.REG_ID,
									NHSOID: i.NHS_OID,
									NHSID:  i.NHS_ID,
								}
								pat.GivenName = i.PIXv3Response.Body.PRPAIN201310UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Given
								pat.FamilyName = i.PIXv3Response.Body.PRPAIN201310UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Family
								for _, pid := range i.PIXv3Response.Body.PRPAIN201310UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.ID {
									switch pid.Root {
									case i.REG_OID:
										pat.REGID = pid.Extension
									case i.NHS_OID:
										pat.NHSID = pid.Extension
									case i.MRN_OID:
										pat.PID = pid.Extension
										pat.PIDOID = i.MRN_OID
									}
								}
								pats := []TUKPatient{}
								pats = append(pats, pat)
								if i.Cache {
									pat_cache[i.Used_PID] = i.Response
									pdq_cache[i.Used_PID] = &pats
								}
							}
						}
					}
				}
			}
		}
	case tukcnst.PDQ_SERVER_TYPE_IHE_PDQV3:
		if tmplt, err = template.New(tukcnst.PDQ_SERVER_TYPE_IHE_PDQV3).Funcs(tukutil.TemplateFuncMap()).Parse(tukcnst.GO_Template_PDQ_V3_Request); err == nil {
			var b bytes.Buffer
			if err = tmplt.Execute(&b, i); err == nil {
				i.Request = b.Bytes()
				if err = i.newIHESOAPRequest(tukcnst.SOAP_ACTION_PDQV3_Request); err == nil {
					if err = xml.Unmarshal(i.Response, &i.PDQv3Response); err == nil {
						if i.PDQv3Response.Body.PRPAIN201306UV02.Acknowledgement.TypeCode.Code != "AA" {
							err = errors.New("acknowledgement code not equal aa, received " + i.PDQv3Response.Body.PRPAIN201306UV02.Acknowledgement.TypeCode.Code)
						} else {
							i.Count, _ = strconv.Atoi(i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.QueryAck.ResultTotalQuantity.Value)
							if i.Count > 0 {
								pat := TUKPatient{
									PIDOID:     i.MRN_OID,
									PID:        i.MRN_ID,
									REGOID:     i.REG_OID,
									REGID:      i.REG_ID,
									NHSOID:     i.NHS_OID,
									NHSID:      i.NHS_ID,
									GivenName:  i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Given,
									FamilyName: i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Family,
									Gender:     i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.AdministrativeGenderCode.Code,
									BirthDate:  i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.BirthTime.Value,
									Street:     i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Addr.StreetAddressLine,
									City:       i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Addr.City,
									State:      i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Addr.State,
									Zip:        i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Addr.PostalCode,
								}
								for _, pid := range i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.ID {
									switch pid.Root {
									case i.REG_OID:
										pat.REGID = pid.Extension
									case i.NHS_OID:
										pat.NHSID = pid.Extension
									case i.MRN_OID:
										pat.PID = pid.Extension
										pat.PIDOID = i.MRN_OID
									}
								}
								pats := []TUKPatient{}
								pats = append(pats, pat)
								if i.Cache {
									pat_cache[i.Used_PID] = i.Response
									pdq_cache[i.Used_PID] = &pats
								}
							}
						}
					}
				}
			}
		}
	case tukcnst.PDQ_SERVER_TYPE_IHE_PIXM:
		i.Request = []byte(i.Server_URL)
		httpReq := tukhttp.PIXmRequest{
			URL:     i.Server_URL,
			PID_OID: i.Used_PID_OID,
			PID:     i.Used_PID,
			Timeout: i.Timeout,
		}
		err = tukhttp.NewRequest(&httpReq)
		i.Response = httpReq.Response
		i.StatusCode = httpReq.StatusCode
		if err == nil {
			if strings.Contains(string(i.Response), "Error") {
				err = errors.New(string(i.Response))
			} else {
				if err := json.Unmarshal(i.Response, &i.PIXmResponse); err == nil {
					log.Printf("%v Patient Entries in Response", i.PIXmResponse.Total)
					i.Count = i.PIXmResponse.Total
					if i.Count > 0 {
						pats := []TUKPatient{}
						for cnt := 0; cnt < len(i.PIXmResponse.Entry); cnt++ {
							rsppat := i.PIXmResponse.Entry[cnt]
							tukpat := TUKPatient{}
							for _, id := range rsppat.Resource.Identifier {
								if id.System == tukcnst.URN_OID_PREFIX+i.REG_OID {
									tukpat.REGID = id.Value
									tukpat.REGOID = i.REG_OID
									log.Printf("Set Reg ID %s %s", tukpat.REGID, tukpat.REGOID)
								}
								if id.Use == "usual" {
									tukpat.PID = id.Value
									tukpat.PIDOID = strings.Split(id.System, ":")[2]
									log.Printf("Set PID %s %s", tukpat.PID, tukpat.PIDOID)
								}
								if id.System == tukcnst.URN_OID_PREFIX+i.NHS_OID {
									tukpat.NHSID = id.Value
									tukpat.NHSOID = i.NHS_OID
									log.Printf("Set NHS ID %s %s", tukpat.NHSID, tukpat.NHSOID)
								}
							}
							gn := ""
							for _, name := range rsppat.Resource.Name {
								for _, n := range name.Given {
									gn = gn + n + " "
								}
							}
							tukpat.GivenName = strings.TrimSuffix(gn, " ")
							tukpat.FamilyName = rsppat.Resource.Name[0].Family
							tukpat.BirthDate = strings.ReplaceAll(rsppat.Resource.BirthDate, "-", "")
							tukpat.Gender = rsppat.Resource.Gender

							if len(rsppat.Resource.Address) > 0 {
								tukpat.Zip = rsppat.Resource.Address[0].PostalCode
								if len(rsppat.Resource.Address[0].Line) > 0 {
									tukpat.Street = rsppat.Resource.Address[0].Line[0]
									if len(rsppat.Resource.Address[0].Line) > 1 {
										tukpat.Town = rsppat.Resource.Address[0].Line[1]
									}
								}
								tukpat.City = rsppat.Resource.Address[0].City
								tukpat.Country = rsppat.Resource.Address[0].Country
							}
							pats = append(pats, tukpat)
							log.Printf("Added Patient %s to response", tukpat.NHSID)
						}
						if i.Cache {
							pat_cache[i.Used_PID] = i.Response
							pdq_cache[i.Used_PID] = &pats
						}
					}
				}
			}
		}
	}
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
func (i *PDQQuery) newIHESOAPRequest(soapaction string) error {
	httpReq := tukhttp.SOAPRequest{
		URL:        i.Server_URL,
		SOAPAction: soapaction,
		Body:       i.Request,
		Timeout:    i.Timeout,
	}
	err := tukhttp.NewRequest(&httpReq)
	i.Response = httpReq.Response
	i.StatusCode = httpReq.StatusCode
	return err
}
