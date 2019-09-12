package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aeternity/aepp-sdk-go/aeternity"
	"github.com/aeternity/aepp-sdk-go/swagguard/node/models"
	rlp "github.com/randomshinichi/rlpae"
	"github.com/spf13/cobra"
)

/* I want to Oraclize this function, associating it with a pubkey.
Actually, I might want to oraclize an entire Struct?
Nah, the function simply distributes to other actual logic. It's okay to just do the function.

so how would I like this to work?

oracle := aeternity.Oraclize(oracleCoreLogic, node1)
oracle.Run()
... running your new Oracle, using node1 to listen and broadcast responses
*/
var (
	oracleQueries = `{"oracle_queries":[{"fee":0,"id":"oq_jUVxdEKipRpedS47QvvNchTxAWfX4SWebhM2dz9zpe4ZkgERZ","oracle_id":"ok_2a1j2Mk9YSmC1gioUq4PWRm3bsv887MbuRVwyv4KaUGoR1eiKi","query":"ov_ZnVjayB5b3UgMEzdCSw=","response":"or_Xfbg4g==","response_ttl":{"type":"delta","value":100},"sender_id":"ak_2a1j2Mk9YSmC1gioUq4PWRm3bsv887MbuRVwyv4KaUGoR1eiKi","sender_nonce":2,"ttl":137},{"fee":0,"id":"oq_2EskxbTcM7nFJNec6hvkj51XqWRHtNySen5HSMvVq4L1t1V95g","oracle_id":"ok_2a1j2Mk9YSmC1gioUq4PWRm3bsv887MbuRVwyv4KaUGoR1eiKi","query":"ov_ZnVjayB5b3UgM/78SrA=","response":"or_Xfbg4g==","response_ttl":{"type":"delta","value":100},"sender_id":"ak_2a1j2Mk9YSmC1gioUq4PWRm3bsv887MbuRVwyv4KaUGoR1eiKi","sender_nonce":5,"ttl":139},{"fee":0,"id":"oq_2NhMjBdKHJYnQjDbAxanmxoXiSiWDoG9bqDgk2MfK2X6AB9Bwx","oracle_id":"ok_2a1j2Mk9YSmC1gioUq4PWRm3bsv887MbuRVwyv4KaUGoR1eiKi","query":"ov_ZnVjayB5b3UgMb3eTQI=","response":"or_Xfbg4g==","response_ttl":{"type":"delta","value":100},"sender_id":"ak_2a1j2Mk9YSmC1gioUq4PWRm3bsv887MbuRVwyv4KaUGoR1eiKi","sender_nonce":3,"ttl":138},{"fee":0,"id":"oq_2oxTwTgMtA4944CPXS4EvovF2Faqh6KvNmEqWEmqymPwWe5dsV","oracle_id":"ok_2a1j2Mk9YSmC1gioUq4PWRm3bsv887MbuRVwyv4KaUGoR1eiKi","query":"ov_ZnVjayB5b3UgMobLbI4=","response":"or_Xfbg4g==","response_ttl":{"type":"delta","value":100},"sender_id":"ak_2a1j2Mk9YSmC1gioUq4PWRm3bsv887MbuRVwyv4KaUGoR1eiKi","sender_nonce":4,"ttl":139}]}`
	i             int
	networkID     = "ae_docker"
	privKey       = "e6a91d633c77cf5771329d3354b3bcef1bc5e032c43d70b6d35af923ce1eb74dcea7ade470c9f99d9d4e400880a86f1d49bb444b62f11a9ebb64bbcfeb73fef3"
	url           = "http://localhost:3013"
	oraclizer     = strings.NewReplacer("ak_", "ok_")
)

func AEUSD() string {
	return "0.22"
}

func choke(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func initialize() (account *aeternity.Account, ctx *aeternity.Context, node *aeternity.Node) {
	account, err := aeternity.AccountFromHexString(privKey)
	if err != nil {
		choke(err)
	}

	ctx, node = aeternity.NewContextFromURL(url, account.Address, false)
	return
}

func sendItOff(tx rlp.Encoder, signingAccount *aeternity.Account, node *aeternity.Node) (err error) {
	_, _, _, blockHeight, _, err := aeternity.SignBroadcastWaitTransaction(tx, signingAccount, node, networkID, 10)
	if err != nil {
		return err
	}
	fmt.Println("Transaction recorded in", blockHeight)
	return nil
}

func oracleInfo(oID string, node *aeternity.Node) (oJSON string, err error) {
	o, err := node.GetOracleByPubkey(oID)
	if err != nil {
		return
	}
	oJ, err := o.MarshalBinary()
	return string(oJ), err
}

var rootCmd = &cobra.Command{
	Use:   "aeoracle",
	Short: "aeoracle is a prototype oracle interface using aepp-sdk-go",
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "gets info on an oracle",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, _, node := initialize()

		ans, err := oracleInfo(args[0], node)
		if err != nil {
			choke(err)
		}
		fmt.Println(ans)
	},
}
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "registers an oracle",
	Run: func(cmd *cobra.Command, args []string) {
		account, ctx, node := initialize()
		registerTx, err := ctx.OracleRegisterTx("just a random question", "it will respond with this", aeternity.Config.Client.Oracles.QueryFee, aeternity.Config.Client.Oracles.QueryTTLType, 500, aeternity.Config.Client.Oracles.VMVersion)
		if err != nil {
			choke(err)
		}

		fmt.Printf("%+v\n", registerTx)
		err = sendItOff(registerTx, account, node)
		if err != nil {
			choke(err)
		}

		oracleID := oraclizer.Replace(account.Address)

		ans, err := oracleInfo(oracleID, node)
		if err != nil {
			choke(err)
		}
		fmt.Println(ans)
	},
}

var extendCmd = &cobra.Command{
	Use:   "extend",
	Short: "extends an oracle's lifetime",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account, ctx, node := initialize()
		extendTx, err := ctx.OracleExtendTx(args[0], aeternity.OracleTTLTypeDelta, 100)
		if err != nil {
			choke(err)
		}

		err = sendItOff(extendTx, account, node)
		if err != nil {
			choke(err)
		}
	},
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "queries an oracle",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account, ctx, node := initialize()
		oracleID := args[0]
		fmt.Println("Spamming the oracle with queries")
		for i := 0; i < 4; i++ {
			tx, err := ctx.OracleQueryTx(oracleID, fmt.Sprintf("fuck you %v", i), aeternity.Config.Client.Oracles.QueryFee, aeternity.Config.Client.Oracles.QueryTTLType, aeternity.Config.Client.Oracles.QueryTTLValue, aeternity.Config.Client.Oracles.ResponseTTLType, aeternity.Config.Client.Oracles.ResponseTTLValue)
			fmt.Printf("%+v", tx)
			if err != nil {
				choke(err)
			}
			err = sendItOff(tx, account, node)
			if err != nil {
				choke(err)
			}
		}
	},
}

var respondCmd = &cobra.Command{
	Use:   "respond",
	Short: "the oracle responds (manually. Use the listen subcommand for automated responding)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account, ctx, node := initialize()
		oracleID := oraclizer.Replace(account.Address)

		fmt.Println("Responding to a Query")
		tx, err := ctx.OracleRespondTx(oracleID, args[0], "My day was fine thank you", aeternity.Config.Client.Oracles.ResponseTTLType, aeternity.Config.Client.Oracles.ResponseTTLValue)
		if err != nil {
			choke(err)
		}

		err = sendItOff(tx, account, node)
		if err != nil {
			choke(err)
		}
	},
}

func getNewOracleQueries(n *aeternity.Node, oID string) (queries []*models.OracleQuery, err error) {
	oq, err := n.GetOracleQueriesByPubkey(oID)
	if err != nil {
		return
	}
	// oq := models.OracleQueries{}
	// oq.UnmarshalBinary([]byte(oracleQueries))

	// Only return new oracle queries
	// fmt.Println("i: ", i, "len(oq)", len(oq.OracleQueries))
	if len(oq.OracleQueries) > i {
		queries = oq.OracleQueries[i+1:]
		i = len(oq.OracleQueries)
		return
	}
	return
}

func listen(node *aeternity.Node, oID string, messages chan []*models.OracleQuery) {
	for {
		q, err := getNewOracleQueries(node, oID)
		// h, err := node.GetHeight()
		if err != nil {
			return
		}

		// fmt.Println("Height:", h, "Just listened", q, err)
		if len(q) > 0 {
			fmt.Printf("Found more oracle queries! %v\n", q)
			messages <- q
		}
		time.Sleep(1 * time.Second)
	}
}

func respond(account *aeternity.Account, ctx *aeternity.Context, node *aeternity.Node, messages chan []*models.OracleQuery) {
	for {
		queries := <-messages
		fmt.Println("Going to respond to []*models.OracleQuery", queries)
		for _, q := range queries {
			tx, err := formulateResponse(q, ctx)
			fmt.Printf("My Response: %+v\n", tx)
			if err != nil {
				fmt.Printf("Error responding to %s: %s", *q.Query, err)
			}
			err = sendItOff(tx, account, node)
			if err != nil {
				fmt.Printf("Error submitting %+v\n", tx)
			}
		}
		fmt.Println("Finished responding to queries")
	}
}

func formulateResponse(q *models.OracleQuery, ctx *aeternity.Context) (r *aeternity.OracleRespondTx, err error) {
	tx, err := ctx.OracleRespondTx(*q.OracleID, *q.ID, "THIS IS AN ANSWER", aeternity.Config.Client.Oracles.ResponseTTLType, aeternity.Config.Client.Oracles.ResponseTTLValue)
	return tx, err
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "oracle listens for incoming queries",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		account, ctx, node := initialize()
		oracleID := oraclizer.Replace(account.Address)
		messages := make(chan []*models.OracleQuery)
		done := make(chan bool)
		go listen(node, oracleID, messages)
		go respond(account, ctx, node, messages)
		<-done
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(extendCmd)
	rootCmd.AddCommand(respondCmd)
	rootCmd.AddCommand(listenCmd)
}

func main() {
	rootCmd.Execute()
}
