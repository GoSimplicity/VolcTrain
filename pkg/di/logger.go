/*
 * Apache License
 * Version 2.0, January 2004
 * http://www.apache.org/licenses/
 *
 * TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION
 *
 * Copyright 2025 Bamboo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package di

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 将日志输出到./logs/cloudops-{日期}.log，并同时输出到控制台
func InitLogger() *zap.Logger {
	// 创建日志目录
	logDir := viper.GetString("log.dir")
	currentTime := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logDir, "cloudops-"+currentTime+".log")

	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("无法创建日志目录")
	}

	// 创建文件输出配置
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,   // 每个日志文件最大10MB就切分
		MaxBackups: 30,   // 保留30个旧文件
		MaxAge:     7,    // 文件最多保存7天
		Compress:   true, // 压缩旧日志文件
		LocalTime:  true, // 使用本地时间
	})

	// 配置日志编码
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 日志等级大写

	// 创建控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 创建 Core
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, zapcore.WarnLevel), // 控制台只输出警告及以上级别
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, zapcore.InfoLevel),       // 文件记录INFO及以上级别
	)

	// 创建 logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}
