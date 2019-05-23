<?php
namespace rpc;

/**
 * Autogenerated by Thrift Compiler (0.12.0)
 *
 * DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
 *  @generated
 */
use Thrift\Base\TBase;
use Thrift\Type\TType;
use Thrift\Type\TMessageType;
use Thrift\Exception\TException;
use Thrift\Exception\TProtocolException;
use Thrift\Protocol\TProtocol;
use Thrift\Protocol\TBinaryProtocolAccelerated;
use Thrift\Exception\TApplicationException;

interface LuckyServiceIf
{
    /**
     * @param int $uid
     * @param string $username
     * @param string $ip
     * @param int $now
     * @param string $app
     * @param string $sign
     * @return \rpc\DataResult
     */
    public function DoLucky($uid, $username, $ip, $now, $app, $sign);
    /**
     * @param int $uid
     * @param string $username
     * @param string $ip
     * @param int $now
     * @param string $app
     * @param string $sign
     * @return \rpc\DataGiftPrize[]
     */
    public function MyPrizeList($uid, $username, $ip, $now, $app, $sign);
}
