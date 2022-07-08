package persistance

import (
	"fmt"
	"strings"
	"time"

	pb "github.com/dislinkt/common/proto/connection_service"
	"github.com/dislinkt/connection_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ConnectionDBStore struct {
	connectionDB *neo4j.Driver
}

func NewConnectionDBStore(client *neo4j.Driver) domain.ConnectionStore {
	return &ConnectionDBStore{
		connectionDB: client,
	}
}

func (store *ConnectionDBStore) CreateConnection(baseUserUuid string, connectUserUuid string) (*pb.NewConnectionResponse, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if checkIfUserExist(baseUserUuid, tx) && checkIfUserExist(connectUserUuid, tx) {

			records, err := tx.Run("MATCH (n:UserNode { uid: $uid}) RETURN n.status", map[string]interface{}{
				"uid": baseUserUuid,
			})
			connectionStatus := ""
			if err != nil {
				return nil, err
			}
			record, err := records.Single()
			if err != nil {
				return nil, err
			}

			status := record.Values[0].(string)
			dateNow := time.Now().Local().Unix()
			fmt.Println(status)
			if status == "PRIVATE" {
				connectionStatus = "REQUEST_SENT"
				_, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $connect_user_uuid  MATCH (u2:UserNode) WHERE u2.uid = $base_user_uuid CREATE (u1)-[r1:CONNECTION {status: $status, date: $date}]->(u2)", map[string]interface{}{
					"connect_user_uuid": connectUserUuid,
					"base_user_uuid":    baseUserUuid,
					"status":            "REQUEST_SENT",
					"date":              dateNow,
				})

				if err != nil {
					return nil, err
				}

			} else {
				connectionStatus = "CONNECTED"
				_, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $connect_user_uuid  MATCH (u2:UserNode) WHERE u2.uid = $base_user_uuid CREATE (u1)-[r2:CONNECTION {status: $status, date: $date}]->(u2) CREATE (u2)-[r1:CONNECTION {status: $status, date: $date}]->(u1)", map[string]interface{}{
					"connect_user_uuid": connectUserUuid,
					"base_user_uuid":    baseUserUuid,
					"status":            "CONNECTED",
					"date":              dateNow,
				})

				if err != nil {
					return nil, err
				}
			}

			fmt.Println(status)

			// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
			return &pb.NewConnectionResponse{
				BaseUserUUID:       baseUserUuid,
				ConnectUserUUID:    connectUserUuid,
				ConnectionResponse: connectionStatus,
			}, nil
		} else {
			return &pb.NewConnectionResponse{
				BaseUserUUID:       baseUserUuid,
				ConnectUserUUID:    connectUserUuid,
				ConnectionResponse: "Connection refused",
			}, nil
		}

	})

	if err != nil {
		return nil, err
	}

	return result.(*pb.NewConnectionResponse), nil
}

func (store *ConnectionDBStore) AcceptConnection(requestSenderUser string, requestApprovalUser string) (*pb.NewConnectionResponse, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if checkIfUserExist(requestSenderUser, tx) && checkIfUserExist(requestApprovalUser, tx) {

			records, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $requestSender MATCH (u2:UserNode) WHERE u2.uid = $requestGet  MATCH (u1)-[r1:CONNECTION]->(u2) return r1.status", map[string]interface{}{
				"requestSender": requestSenderUser,
				"requestGet":    requestApprovalUser,
			})
			connectionStatus := ""
			if err != nil {
				return nil, err
			}
			record, err := records.Single()
			if err != nil {
				return nil, err
			}

			status := record.Values[0].(string)
			dateNow := time.Now().Local().Unix()
			if status == "REQUEST_SENT" {
				connectionStatus = "CONNECTED"
				_, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $requestSender MATCH (u2:UserNode) WHERE u2.uid = $requestGet match  (u1)-[r1:CONNECTION {status: $status0 }]->(u2) set r1.status = $status1 , r1.date = $date  CREATE  (u2)-[r2:CONNECTION {status: $status1, date: $date}]->(u1)", map[string]interface{}{
					"requestSender": requestSenderUser,
					"requestGet":    requestApprovalUser,
					"status0":       "REQUEST_SENT",
					"status1":       "CONNECTED",
					"date":          dateNow,
				})

				if err != nil {
					return nil, err
				}

			} else if status == "CONNECTED" {
				connectionStatus = "CONNECTION EXISTS"
				if err != nil {
					return nil, err
				}
			}

			// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
			return &pb.NewConnectionResponse{
				BaseUserUUID:       requestSenderUser,
				ConnectUserUUID:    requestApprovalUser,
				ConnectionResponse: connectionStatus,
			}, nil
		} else {
			return &pb.NewConnectionResponse{
				BaseUserUUID:       "",
				ConnectUserUUID:    "",
				ConnectionResponse: "Connection refused - user not found",
			}, nil
		}

	})

	if err != nil {
		return nil, err
	}

	return result.(*pb.NewConnectionResponse), nil
}

func (store *ConnectionDBStore) GetAllConnectionForUser(userUid string) (userNodes []*domain.UserNode, error1 error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if !checkIfUserExist(userUid, tx) {
			return &domain.UserNode{
				UserUID: "",
				Status:  "",
			}, nil
		}

		records, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $userUid MATCH (u2:UserNode) WHERE not u2.uid = $userUid match (u2)-[r1:CONNECTION {status: $status}]->(u1) match (u1)-[r2:CONNECTION {status: $status }]->(u2) return u2.uid, u2.status", map[string]interface{}{
			"userUid": userUid,
			"status":  "CONNECTED",
		})

		for records.Next() {
			node := domain.UserNode{UserUID: records.Record().Values[0].(string), Status: domain.Private}
			userNodes = append(userNodes, &node)
		}

		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		return userNodes, nil
	})
	if err != nil {
		return nil, err
	}
	return userNodes, nil
}

func (store *ConnectionDBStore) GetAllConnectionRequestsForUser(userUid string) (userNodes []*domain.UserNode,
	error1 error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		fmt.Println("UUID ")
		fmt.Println(userUid)
		if !checkIfUserExist(userUid, tx) {
			fmt.Println("NE POSTOJI")
			return &domain.UserNode{
				UserUID: "",
				Status:  "",
			}, nil
		}

		records, err := tx.Run("MATCH (u1:UserNode)   MATCH (u2:UserNode) WHERE u2.uid = $userUid match ("+
			"u1)-[r1:CONNECTION {status:$status}]->(u2) return u1.uid, u1.status", map[string]interface{}{
			"userUid": userUid,
			"status":  "REQUEST_SENT",
		})

		for records.Next() {
			fmt.Println(records.Record())
			node := domain.UserNode{UserUID: records.Record().Values[0].(string), Status: domain.Private}
			userNodes = append(userNodes, &node)
			fmt.Println("USAO")
		}

		if err != nil {
			return nil, err
		}
		return userNodes, nil
	})
	if err != nil {
		return nil, err
	}
	return userNodes, nil
}

func (store *ConnectionDBStore) Register(userNode *domain.UserNode) (*domain.UserNode, error) {
	fmt.Println("[ConnectionDBStore Register]")
	fmt.Println(userNode)
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	fmt.Println(session)
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		fmt.Println("linija5")
		if checkIfUserExist(userNode.UserUID, tx) {
			fmt.Println("linija6")
			return &domain.UserNode{
				UserUID: "",
				Status:  "",
			}, nil
		}

		fmt.Println("[ConnectionDBStore Register1]")
		fmt.Println(userNode)
		records, err := tx.Run("CREATE (n:UserNode { uid: $uid, status: $status  }) RETURN n.uid, n.status", map[string]interface{}{
			"uid":    userNode.UserUID,
			"status": userNode.Status,
		})
		fmt.Println("TU SAM")
		if err != nil {
			return nil, err
		}
		record, err := records.Single()
		if err != nil {
			return nil, err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return &domain.UserNode{
			UserUID: record.Values[0].(string),
			Status:  domain.ProfileStatus(record.Values[1].(string)),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*domain.UserNode), nil
}

func checkIfUserExist(uid string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (n:UserNode { uid: $uid }) RETURN n.uid",
		map[string]interface{}{"uid": uid})

	if result != nil && result.Next() && result.Record().Values[0] == uid {
		return true
	}
	return false
}

func (store *ConnectionDBStore) CheckIfUsersConnected(uuid1 string, uuid2 string) (isConnected bool, err error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	res, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		if !checkIfUserExist(uuid1, tx) && !checkIfUserExist(uuid2, tx) {
			return false, nil
		}

		result, _ := tx.Run("match(u1:UserNode) where u1.uid = $userUid "+
			"match (u2:UserNode) where u2.uid = $userUid1 "+
			"MATCH (u1)-[r1:CONNECTION]->(u2) "+
			"RETURN r1.status", map[string]interface{}{
			"userUid":  uuid1,
			"userUid1": uuid2,
		})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return true, err
		} else {
			return false, err
		}

		//re, err := records.Single()
	})

	return res.(bool), err
}

func (store *ConnectionDBStore) UpdateUser(userUUID string, private bool) error {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		fmt.Println("UUID " + userUUID)
		if checkIfUserExist(userUUID, tx) {

			var status string
			if private {
				status = "PRIVATE"
			} else {
				status = "PUBLIC"
			}
			fmt.Println("MENJAM U " + status)

			_, err := tx.Run("MATCH (n:UserNode { uid: $uid}) set n.status = $status",
				map[string]interface{}{
					"uid":    userUUID,
					"status": status,
				})

			if err != nil {
				return nil, err
			}
			return nil, nil
		} else {
			fmt.Println("NEPOSTOJI")
			return nil, nil
		}

	})
	if err != nil {
		return err
	}
	return nil
}

func (store *ConnectionDBStore) BlockUser(currentUser string, blockedUser string) (*pb.BlockedUserStatus, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if checkIfUserExist(currentUser, tx) && checkIfUserExist(blockedUser, tx) {

			records, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $requestSender MATCH (u2:UserNode) WHERE u2.uid = $requestGet  MATCH (u1)-[r1:CONNECTION]->(u2) return r1.status", map[string]interface{}{
				"requestSender": currentUser,
				"requestGet":    blockedUser,
			})

			if err != nil {
				return nil, err
			}
			record, err := records.Single()
			if err != nil {
				return nil, err
			}

			status := record.Values[0].(string)
			dateNow := time.Now().Local().Unix()
			fmt.Println(status)
			if status == "CONNECTED" {
				_, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $currentUser "+
					"Match (u2:UserNode) WHERE u2.uid = $blocked_user "+
					"MATCH (u1)-[r1:CONNECTION {status: 'CONNECTED'}]->(u2) "+
					"MATCH (u2)-[r2:CONNECTION {status: 'CONNECTED'}]->(u1) "+
					"SET r1.status = $status "+
					"SET r2.status = $status1", map[string]interface{}{
					"currentUser":  currentUser,
					"blocked_user": blockedUser,
					"status":       "BLOCK",
					"status1":      "BLOCKED",
					"date":         dateNow,
				})
				fmt.Println("Korisnici izblokirani")

				if err != nil {
					return nil, err
				} else {
					return &pb.BlockedUserStatus{
						CurrentUserUUID:    currentUser,
						BlockedUserUUID:    blockedUser,
						ConnectionResponse: "BLOCK",
					}, nil
				}
			} else if status == "BLOCKED" {
				fmt.Println("uslo")
				_, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $currentUser "+
					"Match (u2:UserNode) WHERE u2.uid = $blocked_user "+
					"MATCH (u1)-[r1:CONNECTION {status: 'BLOCKED'}]->(u2) "+
					"SET r1.status = $status ", map[string]interface{}{
					"currentUser":  currentUser,
					"blocked_user": blockedUser,
					"status":       "BLOCK",
					"date":         dateNow,
				})
				if err != nil {
					return nil, err
				} else {
					return &pb.BlockedUserStatus{
						CurrentUserUUID:    currentUser,
						BlockedUserUUID:    blockedUser,
						ConnectionResponse: "BLOCK",
					}, nil
				}
			} else if status == "BLOCK" {
				return &pb.BlockedUserStatus{
					CurrentUserUUID:    currentUser,
					BlockedUserUUID:    blockedUser,
					ConnectionResponse: "Action refused: user blocked",
				}, nil
			} else {
				return &pb.BlockedUserStatus{
					CurrentUserUUID:    currentUser,
					BlockedUserUUID:    blockedUser,
					ConnectionResponse: "Action refused: users not connected",
				}, nil
			}

		} else {
			return &pb.BlockedUserStatus{
				CurrentUserUUID:    "",
				BlockedUserUUID:    "",
				ConnectionResponse: "Connection refused - user not found",
			}, nil
		}

	})

	if err != nil {
		return nil, err
	}

	return result.(*pb.BlockedUserStatus), nil
}

func (store *ConnectionDBStore) GetAllBlockedForCurrentUser(currentUserUUID string) (blockedUsers []*domain.UserNode, err error) {
	fmt.Println("ConnectionDBStore [GetAllBlockedForCurrentUser]")

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	fmt.Println("ConnectionDBStore [GetAllBlockedForCurrentUser] 1")
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if !checkIfUserExist(currentUserUUID, tx) {
			return &domain.UserNode{
				UserUID: "",
				Status:  "",
			}, nil
		}
		fmt.Println("ConnectionDBStore [GetAllBlockedForCurrentUser] 2")
		records, err := tx.Run("match (u1:UserNode) where u1.uid = $userUid "+
			"match (u2:UserNode) "+
			"MATCH (u1)-[r1:CONNECTION {status: $status}]->(u2) "+
			"return u2.uid, u2.status", map[string]interface{}{
			"userUid": currentUserUUID,
			"status":  "BLOCK",
		})
		fmt.Println(records)
		for records.Next() {
			node := domain.UserNode{}
			if records.Record().Values[1].(string) == "PRIVATE" {
				node = domain.UserNode{UserUID: records.Record().Values[0].(string), Status: domain.Private}
			} else {
				node = domain.UserNode{UserUID: records.Record().Values[0].(string), Status: domain.Public}
			}
			fmt.Println("dodaj korisnika")
			blockedUsers = append(blockedUsers, &node)
		}

		if err != nil {
			return nil, err
		}
		//re, err := records.Single()
		if err != nil {
			return nil, err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return &blockedUsers, nil
	})

	if err != nil {
		return nil, err

	}

	return blockedUsers, nil
}

func (store *ConnectionDBStore) GetAllUserBlockingCurrentUser(currentUserUUID string) (usersThatBlockedYou []*domain.UserNode, err error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if !checkIfUserExist(currentUserUUID, tx) {
			return &domain.UserNode{
				UserUID: "",
				Status:  "",
			}, nil
		}

		records, err := tx.Run("match (u1:UserNode) where u1.uid = $userUid "+
			"match (u2:UserNode) "+
			"MATCH (u2)-[r1:CONNECTION {status: $status}]->(u1) "+
			"return u2.uid, u2.status", map[string]interface{}{
			"userUid": currentUserUUID,
			"status":  "BLOCK",
		})

		for records.Next() {
			node := domain.UserNode{}
			if records.Record().Values[1].(string) == "PRIVATE" {
				node = domain.UserNode{UserUID: records.Record().Values[0].(string), Status: domain.Private}
			} else {
				node = domain.UserNode{UserUID: records.Record().Values[0].(string), Status: domain.Public}
			}

			usersThatBlockedYou = append(usersThatBlockedYou, &node)
		}

		if err != nil {
			return nil, err
		}
		//re, err := records.Single()
		if err != nil {
			return nil, err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return &usersThatBlockedYou, nil
	})

	if err != nil {
		return nil, err

	}

	return usersThatBlockedYou, nil
}

func (store *ConnectionDBStore) RecommendUsersByConnection(currentUserUUID string) (users []*domain.UserNode, err error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		if !checkIfUserExist(currentUserUUID, tx) {
			return &domain.UserNode{
				UserUID: "",
				Status:  "",
			}, nil
		}

		records, err := tx.Run("match(u1:UserNode) where u1.uid = $userUid "+
			"match (u2:UserNode) where u2.uid <> $userUid "+
			"match (u3:UserNode) where u3.uid <> $userUid "+
			"MATCH (u1)-[r1:CONNECTION {status: 'CONNECTED'}]->(u2) "+
			"MATCH (u2)-[r2:CONNECTION {status: 'CONNECTED'}]->(u3) "+
			"RETURN u3.uid, u3.status ", map[string]interface{}{
			"userUid": currentUserUUID,
		})
		var brojac = 0
		for records.Next() {
			result, _ := tx.Run("match(u1:UserNode) where u1.uid = $userUid "+
				"match (u2:UserNode) where u2.uid = $userUid1 "+
				"MATCH (u1)-[r1:CONNECTION]->(u2) "+
				"RETURN r1.status", map[string]interface{}{
				"userUid":  currentUserUUID,
				"userUid1": records.Record().Values[0].(string),
			})

			if !result.Next() {
				brojac = brojac + 1
				if brojac <= 5 {
					fmt.Println("uslo")
					node := domain.UserNode{}
					if records.Record().Values[1].(string) == "PRIVATE" {
						node = domain.UserNode{records.Record().Values[0].(string), domain.Private}
					} else {
						node = domain.UserNode{records.Record().Values[0].(string), domain.Public}
					}
					users = append(users, &node)
				}
			}

		}

		if err != nil {
			return nil, err
		}
		//re, err := records.Single()
		if err != nil {
			return nil, err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return users, nil
	})

	return users, err
}

func (store *ConnectionDBStore) UnblockConnection(currentUser string, blockedUser string) (*pb.BlockedUserStatus, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if checkIfUserExist(currentUser, tx) && checkIfUserExist(blockedUser, tx) {

			records, err := tx.Run("MATCH (u1:UserNode) WHERE u1.uid = $requestSender MATCH (u2:UserNode) WHERE u2.uid = $requestGet  MATCH (u1)-[r1:CONNECTION]->(u2) MATCH (u2)-[r2:CONNECTION]->(u1) return r1.status, r2.status", map[string]interface{}{
				"requestSender": currentUser,
				"requestGet":    blockedUser,
			})

			if err != nil {
				return nil, err
			}
			record, err := records.Single()
			if err != nil {
				return nil, err
			}

			status1 := record.Values[0].(string)
			status2 := record.Values[1].(string)

			if status1 == "BLOCK" && status2 == "BLOCKED" {
				_, err := tx.Run("match (u1:UserNode) where u1.uid = $currentUser "+
					"match (u2:UserNode) where u2.uid = $blockedUser "+
					"MATCH (u1)-[r1:CONNECTION {status: 'BLOCK'}]->(u2) "+
					"MATCH (u2)-[r2:CONNECTION {status: 'BLOCKED'}]->(u1) "+
					"set r1.status = $status, r2.status = $status", map[string]interface{}{
					"currentUser": currentUser,
					"blockedUser": blockedUser,
					"status":      "CONNECTED",
				})
				if err != nil {
					return nil, err
				} else {
					return &pb.BlockedUserStatus{
						CurrentUserUUID:    currentUser,
						BlockedUserUUID:    blockedUser,
						ConnectionResponse: "Connection status - CONNECTED",
					}, nil
				}
			} else if status1 == "BLOCK" && status2 == "BLOCK" {

				_, err := tx.Run("match (u1:UserNode) where u1.uid = $currentUser "+
					"match (u2:UserNode) where u2.uid = $blockedUser "+
					"MATCH (u1)-[r1:CONNECTION {status: 'BLOCK'}]->(u2) "+
					"MATCH (u2)-[r2:CONNECTION {status: 'BLOCK'}]->(u1) "+
					"set r1.status = $status", map[string]interface{}{
					"currentUser": currentUser,
					"blockedUser": blockedUser,
					"status":      "BLOCKED",
				})
				if err != nil {
					return nil, err
				} else {
					return &pb.BlockedUserStatus{
						CurrentUserUUID:    currentUser,
						BlockedUserUUID:    blockedUser,
						ConnectionResponse: "Connection status - BLOCKED",
					}, nil
				}
			} else {
				return &pb.BlockedUserStatus{
					CurrentUserUUID:    "",
					BlockedUserUUID:    "",
					ConnectionResponse: "Connection refused - user not found",
				}, nil
			}
		} else {
			return &pb.BlockedUserStatus{
				CurrentUserUUID:    "",
				BlockedUserUUID:    "",
				ConnectionResponse: "Connection refused - user not found",
			}, nil
		}

	})

	if err != nil {
		return nil, err
	}

	return result.(*pb.BlockedUserStatus), nil
}

func (store *ConnectionDBStore) InsertField(name string) (string, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		records, err := tx.Run("CREATE (n:Field { name: $name  }) RETURN n.name", map[string]interface{}{
			"name": name,
		})

		if err != nil {
			return "", err
		}
		record, err := records.Single()
		if err != nil {
			return "", err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return record.Values[0].(string), nil

	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func (store *ConnectionDBStore) InsertSkill(name string) (string, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	fmt.Println("usloo")

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		records, err := tx.Run("CREATE (n:Skill { name: $name  }) RETURN n.name", map[string]interface{}{
			"name": name,
		})

		if err != nil {
			return "", err
		}
		record, err := records.Single()
		if err != nil {
			return "", err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return record.Values[0].(string), nil

	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func (store *ConnectionDBStore) InsertJobOffer(jobOffer domain.JobOffer) (string, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	fmt.Println("ConnectionDBStore: InsertJobOffer")
	fmt.Println(jobOffer)

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		dateNow := time.Now().Local().Unix()
		records, err := tx.Run("CREATE (n:JobOffer { Id: $id,Position: $position,DatePosted: $datePosted, Duration: $duration, "+
			"Location: $location, Title: $title,Field: $field,Preconditions: $preconditions}) RETURN n.Id", map[string]interface{}{
			"id":            jobOffer.Id,
			"position":      jobOffer.Position,
			"datePosted":    dateNow,
			"duration":      jobOffer.Duration,
			"location":      jobOffer.Location,
			"title":         jobOffer.Title,
			"field":         jobOffer.Field,
			"preconditions": jobOffer.Preconditions,
		})

		if err != nil {
			return "", err
		}

		records, err = tx.Run("match (j:JobOffer) where j.Id = $id "+
			"match (f:Field) where f.name = $field "+
			"CREATE (j)-[r1:FIELDS]->(f)", map[string]interface{}{
			"id":    jobOffer.Id,
			"field": jobOffer.Field,
		})

		s_array := strings.Split(jobOffer.Preconditions, ",")

		for _, s := range s_array {

			records, err = tx.Run("match (j:JobOffer) where j.Id = $id "+
				"match (s:Skill) where s.name = $skill "+
				"CREATE (j)-[r1:SKILLS]->(s) return j.Id", map[string]interface{}{
				"id":    jobOffer.Id,
				"skill": s,
			})
		}

		record, err := records.Single()
		if err != nil {
			return "", err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return record.Values[0].(string), nil

	})

	if err != nil {
		return "", err
	}

	fmt.Println("ConnectionDBStore: InsertJobOffer ", result.(string))

	return result.(string), err
}

func (store *ConnectionDBStore) InsertSkillToUser(name string, uuid string) (string, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		records, err := tx.Run("match (u1:UserNode) where u1.uid = $uuid "+
			"match (s:Skill) where s.name = $name "+
			"CREATE (u1)-[r1:SKILLS]->(s) return u1.uid", map[string]interface{}{
			"name": name,
			"uuid": uuid,
		})

		if err != nil {
			return "", err
		}
		record, err := records.Single()
		if err != nil {
			return "", err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return record.Values[0].(string), nil

	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func (store *ConnectionDBStore) InsertFieldToUser(name string, uuid string) (string, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		records, err := tx.Run("match (u1:UserNode) where u1.uid = $uuid "+
			"match (s:Field) where s.name = $name "+
			"CREATE (u1)-[r1:FIELDS]->(s) return u1.uid", map[string]interface{}{
			"name": name,
			"uuid": uuid,
		})

		if err != nil {
			return "", err
		}
		record, err := records.Single()
		if err != nil {
			return "", err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return record.Values[0].(string), nil

	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func (store *ConnectionDBStore) RecommendJobBySkill(userUid string) (jobs []*domain.JobOffer, err error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if !checkIfUserExist(userUid, tx) {
			return nil, nil
		}

		records, err := tx.Run("match (u1:UserNode) where u1.uid = $uuid "+
			"match (s:Skill) "+
			"match (j:JobOffer) "+
			"MATCH (u1)-[r1:SKILLS]->(s) "+
			"MATCH (j)-[r2:SKILLS]->(s) "+
			"return DISTINCT j.Id , j.Position, j.Preconditions,j.DatePosted, j.Duration, j.Location, j.Title, j.Field "+
			"LIMIT 5", map[string]interface{}{
			"uuid": userUid,
		})

		for records.Next() {

			node := domain.JobOffer{
				Id:            records.Record().Values[0].(string),
				Position:      records.Record().Values[1].(string),
				Preconditions: records.Record().Values[2].(string),
				Duration:      records.Record().Values[4].(string),
				Location:      records.Record().Values[5].(string),
				Title:         records.Record().Values[6].(string),
				Field:         records.Record().Values[7].(string),
			}

			jobs = append(jobs, &node)
		}

		if err != nil {
			return nil, err
		}
		//re, err := records.Single()
		if err != nil {
			return nil, err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return jobs, nil
	})

	return jobs, err
}

func (store *ConnectionDBStore) RecommendJobByField(userUid string) (jobs []*domain.JobOffer, err error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		if !checkIfUserExist(userUid, tx) {
			return nil, nil
		}

		records, err := tx.Run("match (u1:UserNode) where u1.uid = $uuid "+
			"match (s:Field) "+
			"match (j:JobOffer) "+
			"MATCH (u1)-[r1:FIELDS]->(s) "+
			"MATCH (j)-[r2:FIELDS]->(s) "+
			"return DISTINCT j.Id , j.Position, j.Preconditions,j.DatePosted, j.Duration, j.Location, j.Title, j.Field "+
			"LIMIT 5", map[string]interface{}{
			"uuid": userUid,
		})

		for records.Next() {

			node := domain.JobOffer{
				Id:            records.Record().Values[0].(string),
				Position:      records.Record().Values[1].(string),
				Preconditions: records.Record().Values[2].(string),
				Duration:      records.Record().Values[4].(string),
				Location:      records.Record().Values[5].(string),
				Title:         records.Record().Values[6].(string),
				Field:         records.Record().Values[7].(string),
			}

			jobs = append(jobs, &node)
		}

		if err != nil {
			return nil, err
		}
		//re, err := records.Single()
		if err != nil {
			return nil, err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return jobs, nil
	})

	return jobs, err
}

func (store *ConnectionDBStore) DeleteAllSkills() (res string, err error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		_, err := tx.Run("match (s:Skill) "+
			"match (u1:UserNode) "+
			"match (j:JobOffer) "+
			"match (u1)-[r1:SKILLS]->(s) "+
			"match (j)-[r2:SKILLS]->(s) "+
			"delete r1,r2,s", map[string]interface{}{})

		if err != nil {
			return nil, err
		}
		//re, err := records.Single()
		if err != nil {
			return nil, err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return "done", nil
	})

	return "done", err
}

func (store *ConnectionDBStore) DeleteAllFields() (res string, err error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		_, err := tx.Run("match (s:Field) "+
			"match (u1:UserNode) "+
			"match (j:JobOffer) "+
			"match (u1)-[r1:FIELDS]->(s) "+
			"match (j)-[r2:FIELDS]->(s) "+
			"delete r1,r2,s", map[string]interface{}{})

		if err != nil {
			return nil, err
		}
		//re, err := records.Single()
		if err != nil {
			return nil, err
		}
		// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
		return "done", nil
	})

	return "done", err
}

func (store *ConnectionDBStore) DeleteSkillForUser(userUUID string, skillName string) (res string, err error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		_, err := tx.Run("match (s:Skill) where s.name = $name "+
			"match (u1:UserNode) where u1.uid = $uuid "+
			"match (u1)-[r1:SKILLS]->(s) "+
			"delete r1", map[string]interface{}{
			"uuid": userUUID,
			"name": skillName,
		})

		if err != nil {
			return nil, err
		}
		return "done", nil
	})

	return "done", err
}

func (store *ConnectionDBStore) DeleteFieldForUser(userUUID string, fieldName string) (res string, err error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		_, err := tx.Run("match (s:Field) where s.name = $name "+
			"match (u1:UserNode) where u1.uid = $uuid "+
			"match (u1)-[r1:FIELDS]->(s) "+
			"delete r1", map[string]interface{}{
			"uuid": userUUID,
			"name": fieldName,
		})

		if err != nil {
			return nil, err
		}
		return "done", nil
	})

	return "done", err
}

func (store *ConnectionDBStore) UpdateSkillForUser(userUUID string, skillNameOld string, skillNameNew string) (res string, err error) {

	res, err = store.DeleteSkillForUser(userUUID, skillNameOld)
	if res == "done" {
		res, err = store.InsertSkillToUser(userUUID, skillNameNew)
		if res != "done" {
			return "error", err
		}
	}

	return "done", err
}

func (store *ConnectionDBStore) UpdateFieldForUser(userUUID string, fieldNameOld string, fieldNameNew string) (res string, err error) {

	res, err = store.DeleteFieldForUser(userUUID, fieldNameOld)
	if res == "done" {
		res, err = store.InsertFieldToUser(userUUID, fieldNameNew)
		if res != "done" {
			return "error", err
		}
	}

	return "done", err
}
