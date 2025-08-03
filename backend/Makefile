# VolcTrain Backend 测试 Makefile
# 使用方法: make test-unit, make test-integration, make test-all

# 变量定义
GO := go
GOTEST := $(GO) test
GOBUILD := $(GO) build
GOMOD := $(GO) mod
GOGENERATE := $(GO) generate

# 测试相关变量
TEST_TIMEOUT := 30m
TEST_VERBOSE := -v
TEST_RACE := -race
TEST_COVERAGE := -coverprofile=coverage.out
TEST_COVERAGE_HTML := -html=coverage.out -o coverage.html

# 数据库配置
TEST_DB_HOST := localhost
TEST_DB_NAME := volctrain_test
TEST_DB_USER := volctrain
TEST_DB_PASS := volctrain

# Redis配置
TEST_REDIS_HOST := localhost
TEST_REDIS_PORT := 6379

# 目录定义
SRC_DIR := .
TEST_DIR := ./test
UNIT_TEST_DIR := $(TEST_DIR)/unit
INTEGRATION_TEST_DIR := $(TEST_DIR)/integration
MOCK_DIR := $(TEST_DIR)/mocks

.PHONY: help test-deps test-unit test-integration test-performance test-security test-all coverage clean

# 默认目标
help:
	@echo "VolcTrain Backend 测试命令"
	@echo ""
	@echo "测试依赖管理:"
	@echo "  test-deps        安装测试依赖"
	@echo "  test-setup       设置测试环境"
	@echo "  test-cleanup     清理测试环境"
	@echo ""
	@echo "单元测试:"
	@echo "  test-unit        运行所有单元测试"
	@echo "  test-unit-model  运行模型层单元测试"
	@echo "  test-unit-logic  运行业务逻辑层单元测试"
	@echo "  test-unit-pkg    运行工具包单元测试"
	@echo ""
	@echo "集成测试:"
	@echo "  test-integration 运行所有集成测试"
	@echo "  test-api         运行API集成测试"
	@echo "  test-db          运行数据库集成测试"
	@echo ""
	@echo "性能和安全测试:"
	@echo "  test-performance 运行性能测试"
	@echo "  test-security    运行安全测试"
	@echo "  test-load        运行负载测试"
	@echo ""
	@echo "综合测试:"
	@echo "  test-all         运行所有测试"
	@echo "  test-ci          CI/CD模式运行测试"
	@echo ""
	@echo "测试报告:"
	@echo "  coverage         生成测试覆盖率报告"
	@echo "  coverage-html    生成HTML覆盖率报告"
	@echo "  test-report      生成完整测试报告"

# 安装测试依赖
test-deps:
	@echo "安装测试依赖..."
	$(GOMOD) tidy
	$(GOMOD) download
	@echo "依赖安装完成"

# 生成Mock文件
generate-mocks:
	@echo "生成Mock文件..."
	$(GOGENERATE) $(MOCK_DIR)/...
	@echo "Mock文件生成完成"

# 设置测试环境
test-setup:
	@echo "设置测试环境..."
	@echo "检查数据库连接..."
	@mysql -h$(TEST_DB_HOST) -u$(TEST_DB_USER) -p$(TEST_DB_PASS) -e "CREATE DATABASE IF NOT EXISTS $(TEST_DB_NAME);" || true
	@echo "初始化测试数据库..."
	@mysql -h$(TEST_DB_HOST) -u$(TEST_DB_USER) -p$(TEST_DB_PASS) $(TEST_DB_NAME) < sql/00_create_schema.sql || true
	@echo "测试环境设置完成"

# 清理测试环境
test-cleanup:
	@echo "清理测试环境..."
	@rm -f coverage.out coverage.html
	@rm -f test-results.xml
	@rm -rf test-reports/
	@echo "清理完成"

# 运行单元测试
test-unit:
	@echo "运行单元测试..."
	TEST_DB_HOST=$(TEST_DB_HOST) TEST_DB_NAME=$(TEST_DB_NAME) TEST_DB_USER=$(TEST_DB_USER) TEST_DB_PASS=$(TEST_DB_PASS) \
	$(GOTEST) $(TEST_VERBOSE) $(TEST_RACE) $(TEST_COVERAGE) -timeout $(TEST_TIMEOUT) \
		./internal/logic/... \
		./model/... \
		./pkg/... \
		-tags=unit
	@echo "单元测试完成"

# 运行模型层单元测试
test-unit-model:
	@echo "运行模型层单元测试..."
	TEST_DB_HOST=$(TEST_DB_HOST) TEST_DB_NAME=$(TEST_DB_NAME) TEST_DB_USER=$(TEST_DB_USER) TEST_DB_PASS=$(TEST_DB_PASS) \
	$(GOTEST) $(TEST_VERBOSE) $(TEST_RACE) -timeout $(TEST_TIMEOUT) ./model/... -tags=unit

# 运行业务逻辑层单元测试
test-unit-logic:
	@echo "运行业务逻辑层单元测试..."
	TEST_DB_HOST=$(TEST_DB_HOST) TEST_DB_NAME=$(TEST_DB_NAME) TEST_DB_USER=$(TEST_DB_USER) TEST_DB_PASS=$(TEST_DB_PASS) \
	$(GOTEST) $(TEST_VERBOSE) $(TEST_RACE) -timeout $(TEST_TIMEOUT) ./internal/logic/... -tags=unit

# 运行工具包单元测试
test-unit-pkg:
	@echo "运行工具包单元测试..."
	$(GOTEST) $(TEST_VERBOSE) $(TEST_RACE) -timeout $(TEST_TIMEOUT) ./pkg/... -tags=unit

# 运行集成测试
test-integration:
	@echo "运行集成测试..."
	TEST_DB_HOST=$(TEST_DB_HOST) TEST_DB_NAME=$(TEST_DB_NAME) TEST_DB_USER=$(TEST_DB_USER) TEST_DB_PASS=$(TEST_DB_PASS) \
	TEST_REDIS_HOST=$(TEST_REDIS_HOST) TEST_REDIS_PORT=$(TEST_REDIS_PORT) \
	$(GOTEST) $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./tests/integration/... -tags=integration
	@echo "集成测试完成"

# 运行API集成测试
test-api:
	@echo "运行API集成测试..."
	TEST_DB_HOST=$(TEST_DB_HOST) TEST_DB_NAME=$(TEST_DB_NAME) TEST_DB_USER=$(TEST_DB_USER) TEST_DB_PASS=$(TEST_DB_PASS) \
	$(GOTEST) $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./tests/integration/api/... -tags=integration

# 运行数据库集成测试
test-db:
	@echo "运行数据库集成测试..."
	TEST_DB_HOST=$(TEST_DB_HOST) TEST_DB_NAME=$(TEST_DB_NAME) TEST_DB_USER=$(TEST_DB_USER) TEST_DB_PASS=$(TEST_DB_PASS) \
	$(GOTEST) $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./tests/integration/database/... -tags=integration

# 运行性能测试
test-performance:
	@echo "运行性能测试..."
	TEST_DB_HOST=$(TEST_DB_HOST) TEST_DB_NAME=$(TEST_DB_NAME) TEST_DB_USER=$(TEST_DB_USER) TEST_DB_PASS=$(TEST_DB_PASS) \
	$(GOTEST) $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) -bench=. -benchmem ./tests/performance/... -tags=performance
	@echo "性能测试完成"

# 运行安全测试
test-security:
	@echo "运行安全测试..."
	$(GOTEST) $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./tests/security/... -tags=security
	@echo "安全测试完成"

# 运行负载测试
test-load:
	@echo "运行负载测试..."
	@echo "启动vegeta负载测试..."
	@if command -v vegeta >/dev/null 2>&1; then \
		cd scripts && ./performance-test.sh; \
	else \
		echo "vegeta未安装，跳过负载测试"; \
	fi

# 运行所有测试
test-all: test-unit test-integration test-performance test-security
	@echo "所有测试完成"

# CI/CD模式测试
test-ci:
	@echo "CI/CD模式测试..."
	TEST_DB_HOST=$(TEST_DB_HOST) TEST_DB_NAME=$(TEST_DB_NAME) TEST_DB_USER=$(TEST_DB_USER) TEST_DB_PASS=$(TEST_DB_PASS) \
	$(GOTEST) $(TEST_RACE) $(TEST_COVERAGE) -timeout $(TEST_TIMEOUT) -json ./... > test-results.json
	@echo "CI/CD测试完成"

# 生成测试覆盖率报告
coverage: test-unit
	@echo "生成测试覆盖率报告..."
	$(GO) tool cover -func=coverage.out
	@echo "覆盖率统计:"
	@$(GO) tool cover -func=coverage.out | grep "total:" | awk '{print "总覆盖率: " $$3}'

# 生成HTML覆盖率报告
coverage-html: test-unit
	@echo "生成HTML覆盖率报告..."
	$(GO) tool cover $(TEST_COVERAGE_HTML)
	@echo "HTML报告已生成: coverage.html"

# 检查测试覆盖率
coverage-check: coverage
	@echo "检查测试覆盖率..."
	@COVERAGE=$$($(GO) tool cover -func=coverage.out | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
	if [ "$${COVERAGE%.*}" -lt 80 ]; then \
		echo "警告: 测试覆盖率 $${COVERAGE}% 低于80%"; \
		exit 1; \
	else \
		echo "测试覆盖率 $${COVERAGE}% 达标"; \
	fi

# 生成完整测试报告
test-report: test-all coverage-html
	@echo "生成完整测试报告..."
	@mkdir -p test-reports
	@cp coverage.html test-reports/
	@$(GO) test -json ./... > test-reports/test-results.json
	@echo "测试报告已生成在 test-reports/ 目录"

# 代码静态分析
lint:
	@echo "运行代码静态分析..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint未安装，跳过静态分析"; \
	fi

# 代码格式化
fmt:
	@echo "格式化代码..."
	$(GO) fmt ./...

# 代码整理
tidy:
	@echo "整理依赖..."
	$(GOMOD) tidy

# 构建测试
build-test:
	@echo "测试构建..."
	$(GOBUILD) -o /dev/null ./...

# 验证所有检查
verify: fmt tidy lint build-test test-all coverage-check
	@echo "所有验证检查完成"

# 快速测试（仅运行重要测试）
test-quick:
	@echo "快速测试..."
	$(GOTEST) $(TEST_VERBOSE) -short ./internal/logic/... ./model/...

# 并行测试
test-parallel:
	@echo "并行测试..."
	$(GOTEST) $(TEST_VERBOSE) -parallel 4 ./...

# 基准测试
benchmark:
	@echo "运行基准测试..."
	$(GOTEST) -bench=. -benchmem ./...

# 竞争检测
test-race:
	@echo "竞争检测测试..."
	$(GOTEST) $(TEST_RACE) ./...

# 内存泄漏检测
test-memory:
	@echo "内存泄漏检测..."
	$(GOTEST) -memprofile=mem.prof ./...

# 清理所有生成文件
clean: test-cleanup
	@echo "清理所有生成文件..."
	@rm -f *.prof
	@rm -f *.test
	@rm -rf vendor/