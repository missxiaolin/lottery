<?php
// 加载基础的类库代码
error_reporting(E_ALL);
require_once __DIR__.'/lib/Thrift/ClassLoader/ThriftClassLoader.php';
use Thrift\ClassLoader\ThriftClassLoader;

$GEN_DIR = realpath(dirname(__FILE__)).'/../../';
include $GEN_DIR . "Types.php";
include $GEN_DIR . "LuckyService.php";

$loader = new ThriftClassLoader();
$loader->registerNamespace('Thrift',__DIR__.'/lib');
$loader->registerDefinition('rpc',$GEN_DIR);
$loader->register();

use Thrift\Protocol\TBinaryProtocol;
use \Thrift\Protocol\TJSONProtocol;
use Thrift\Transport\TSocket;
use Thrift\Transport\THttpClient;
use Thrift\Transport\TBufferedTransport;
use Thrift\Exception\TException;

// 得到thrift客户端连接对象
function getClient() {
    try {
        $socket = new THttpClient('localhost',8080, '/rpc', 'http');

        $transport = new TBufferedTransport($socket,1024,1024);
        $protocol = new TJSONProtocol($transport);
        $client = new \rpc\LuckyServiceClient($protocol);
        $transport->open();

        return $client;
    } catch (\Exception $e) {
        print 'TException:'.$e->getMessage().PHP_EOL;
    }
//        $transport->close();
}

$client = getClient();

$uid = 20;
$username = "admin";
$ip = "127.0.0.1";
$now = time();
$app = "web";
$sign = md5("0123456789abcdefuid=$uid&username=$username&ip=$ip&now=$now&app=$app");

try {
    echo "DoLucky\n";
    $rs = $client->DoLucky($uid, $username, $ip, $now, $app, $sign);
    var_dump($rs);
} catch (Exception $err) {
    print_r($err);
}

try {
    echo "MyPrizeList\n";
    $rs = $client->MyPrizeList($uid, $username, $ip, $now, $app, $sign);
    var_dump($rs);
} catch (Exception $err) {
    print_r($err);
}
