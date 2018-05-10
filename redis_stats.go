/*
The purpose of this file is to reproduce the following PHP code with equivalent Go. The output
is written to a text file and read by an `item` in Zabbix.
*/

/*

foreach ($redis->info('COMMANDSTATS') as $cmd => $data) {
    if (!preg_match('!^calls=([0-9]+),usec=([0-9]+)!', $data, $m)) {
        continue;
    }

    $cmd = str_replace('cmdstat_', '', $cmd);
    switch ($cmd) {
        case 'sscan':
        case 'select':
        case 'zrevrangebyscore':
        case 'exists':
        case 'getbit':
        case 'get':
        case 'mget':
        case 'scan':
        case 'scard':
        case 'zcount':
        case 'zrangebyscore':
        case 'sismember':
        case 'strlen':
        case 'type':
        case 'smembers':
        case 'zcard':
        case 'hget':
        case 'hlen':
        case 'hkeys':
            $out['read_ops'] += $m[1];
            $out['read_usec'] += $m[2];
            break;

        case 'zadd':
        case 'setnx':
        case 'zremrangebyscore':
        case 'del':
        case 'rpush':
        case 'zrem':
        case 'lpop':
        case 'getset':
        case 'incr':
        case 'setbit':
        case 'expireat':
        case 'setex':
        case 'expire':
        case 'srem':
        case 'sadd':
        case 'set':
        case 'hsetnx':
        case 'hdel':
        case 'hset':
            $out['write_ops'] += $m[1];
            $out['write_usec'] += $m[2];
            break;

        default:
            $out['other_ops'] += $m[1];
            $out['other_usec'] += $m[2];
            break;
    }
    $out['total_ops'] += $m[1];
    $out['total_usec'] += $m[2];
}

file_put_contents(STATE_FILE, json_encode($out));

if (!$prev) {
    return;
}

$txt = [];
foreach ($out as $k => $v) {
    if (in_array($k, ['used_memory', 'connected_clients', 'instantaneous_input_kbps', 'instantaneous_output_kbps'])) {
        $txt[$k] = $v;
        continue;
    }
    $txt[$k] = ($v - $prev[$k]);
}

$hits = $out['keyspace_hits'] - $prev['keyspace_hits'];
$miss = $out['keyspace_misses'] - $prev['keyspace_misses'];

$txt['hit_ratio'] = $txt['keyspace_hits'] ? round($txt['keyspace_hits'] / ($txt['keyspace_hits'] + $txt['keyspace_misses']) * 100, 2) : 0;

foreach (['read', 'write', 'other', 'total'] as $v) {
    $txt["{$v}_latency"] = $txt["{$v}_ops"] ? round($txt["{$v}_usec"] / $txt["{$v}_ops"] / $txt['uptime_in_seconds'], 2) : 0;
    $txt["{$v}_ops"] = round($txt["{$v}_ops"] / $txt['uptime_in_seconds']);
}

$o = '';
foreach ($txt as $k => $v) {
    $o .= "{$k} {$v}\n";
}

file_put_contents(OUT_FILE, $o);
*/

package main 

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "io/ioutil"
    "log"
    "github.com/gomodule/redigo/redis"
)

var (
    // flag vars
    host    string 
    port    string 
)

type CoreStats struct {
    used_memory int
    connected_clients int
    instantaneous_input_kbps int
    instantaneous_output_kbps int
    keyspace_hits int
    keyspace_misses int
    expired_keys int
    evicted_keys int
    read_ops int
    write_ops int
    other_ops int
    total_ops int
    read_usec int
    write_usec int
    other_usec int
    total_usec int
}

func init() {
    flag.StringVar(&host, "host", "localhost", "Redis hostname")
    flag.StringVar(&port, "port", "6379", "Redis port")
}

func main() {

    flag.Parse()

    // Confirm that redis can be connected to. Error out, if not.
    var pool = newPool()
    c := pool.Get()
    defer c.Close()

    _, err := c.Do("PING")
    if err != nil {
      log.Fatal("Can't connect to the Redis database")
    }

    // Check is JSON state file exists and generate if it doesn't
    statefile := "/home/jason/redis_statefile.txt"
    if _, err := os.Stat(statefile); os.IsNotExist(err) {
        var jsonBlob = []byte(`
            {"uptime_in_seconds":0,
            "used_memory":0,
            "connected_clients":0,
            "instantaneous_input_kbps":0,
            "instantaneous_output_kbps":0,
            "keyspace_misses":0,
            "keyspace_hits":0,
            "expired_keys":0,
            "evicted_keys":0,
            "read_ops":0,
            "write_ops":0,
            "other_ops":0,
            "total_ops":0,
            "read_usec":0,
            "write_usec":0,
            "other_usec":0,
            "total_usec":0
            }`)

    /* Some of the above values can be filled by redis-cli commands
        uptime_in_seconds
        used_memory
        connected_clients
        instantaneous_input_kbps
        instantaneous_output_kbps
        keyspace_misses
        keyspace_hits
        expired_keys
        evicted_keys
    */

        corestats := CoreStats{}
        err := json.Unmarshal(jsonBlob, &corestats)
        if err != nil {
            fmt.Println(err)
            return
        }

        // The writing to file is not working properly. Looks like corestatsJson isn't correct. 
        // Will debug later on.
        corestatsJson, _ := json.Marshal(corestats)
        err = ioutil.WriteFile("/home/jason/redis_statefile.txt", corestatsJson, 0644)
        fmt.Printf("%+v", corestats)

    } else {
        fmt.Printf("File exists.\n")
    }

}

func newPool() *redis.Pool {
    return &redis.Pool{
        MaxIdle:   80,
        MaxActive: 1000, // max number of connections
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", host+":"+port)
            if err != nil {
                panic(err.Error())
            }
            return c, err
        },
    }
}