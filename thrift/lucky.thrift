namespace go rpc
namespace php rpc
namespace java rpc

# 奖品详情
struct DataGiftPrize{
    1: i64    Id             = 0
    2: string Title          = ""
    3: string Img            = ""
    4: i64    Displayorder   = 0
    5: i64    Gtype          = 0
    6: string Gdata          = ""
}
# 返回值
struct DataResult{
    1:i64           Code
    2:string        Msg
    3:DataGiftPrize Gift
}

# 服务接口
service LuckyService {
    # 抽奖的方法
    DataResult DoLucky(1:i64 uid, 2:string username, 3:string ip, 4:i64 now, 5:string app, 6:string sign),
    list<DataGiftPrize> MyPrizeList(1:i64 uid, 2:string username, 3:string ip, 4:i64 now, 5:string app, 6:string sign),
}
