package main

import (
	"errors"
	"time"
)

var (
	start_timestamp int64 = 1557489395327                           // 开始时间戳
	sequence_bit    int64 = 12                                      // 序列号占用位数
	machine_int     int64 = 10                                      // 机器标识所占位数
	timestamp_left  int64 = sequence_bit + machine_int              // 时间戳位移位数
	max_sequence    int64 = -1 ^ (-1 << sequence_bit)               // 最大序列号
	max_machine_id  int64 = -1 ^ (-1 << machine_int)                // 最大机器编号
	machineIdPart   int64 = (9123 & max_machine_id) << sequence_bit // 生成id 机器标识部分
	sequence        int64 = 0                                       // 序列号
	last_stamp      int64 = -1
)

func NextId() (int64, error) {
	currentStamp := time.Now().UnixNano() / 1e6 // 当前是时间戳（毫秒）
	// 当前时间小于最后生成的时间
	if currentStamp < last_stamp {
		err := errors.New("时钟已经回拨")
		return 0, err
	}
	// 当前时间等于最后生成时间，阻塞获取下一毫秒
	if currentStamp == last_stamp {
		sequence = (sequence + 1) & max_sequence
		if sequence == 0 {
			currentStamp = getNextMill()
		}
	} else {
		sequence = 0
	}
	// 修改最后生成时间
	last_stamp = currentStamp

	return (currentStamp-start_timestamp)<<timestamp_left | machineIdPart | sequence, nil
}

func getNextMill() int64 {
	mill := time.Now().UnixNano() / 1e6
	for mill <= last_stamp {
		mill = time.Now().UnixNano() / 1e6
	}
	return mill
}
