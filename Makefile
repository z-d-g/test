BINARY  := md-cli
MODULE  := ./cmd/md-cli
CGO     := CGO_ENABLED=0

# Build flags
LDFLAGS := -s -w
GOFLAGS := -trimpath

# Version metadata (override via env or make args)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

EXTLDFLAGS := $(LDFLAGS) -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Run target file
FILE    ?= README.md

# Installation prefix
PREFIX  ?= $(HOME)/.local
BINDIR  := $(PREFIX)/bin

# Platform detection
OS      := $(shell uname -s | tr A-Z a-z)
ARCH    := $(shell uname -m)

# Normalize arch names
ifeq ($(ARCH),x86_64)
  ARCH := amd64
endif
ifeq ($(ARCH),aarch64)
  ARCH := arm64
endif

# ── Default ────────────────────────────────────────────────────────
.PHONY: all
all: build

# ── Local build ────────────────────────────────────────────────────
.PHONY: build
build:
	$(CGO) go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o bin/$(BINARY) $(MODULE)

# ── Build with version info ────────────────────────────────────────
.PHONY: release
release:
	$(CGO) go build $(GOFLAGS) -ldflags="$(EXTLDFLAGS)" -o bin/$(BINARY) $(MODULE)

# ── Cross-compile all platforms ────────────────────────────────────
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: build-all
build-all:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		ext=""; \
		[ "$$os" = "windows" ] && ext=".exe"; \
		output="dist/$(BINARY)-$$os-$$arch$$ext"; \
		echo "  $$output"; \
		GOOS=$$os GOARCH=$$arch $(CGO) go build $(GOFLAGS) \
			-ldflags="$(EXTLDFLAGS)" \
			-o "$$output" $(MODULE) || exit 1; \
	done
	@echo "✓ All binaries in dist/"

# ── Create release tarballs/zip ────────────────────────────────────
.PHONY: dist
dist: build-all
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		ext=""; \
		[ "$$os" = "windows" ] && ext=".exe"; \
		src="dist/$(BINARY)-$$os-$$arch$$ext"; \
		if [ "$$os" = "windows" ]; then \
			(cd dist && zip -q "$(BINARY)-$$os-$$arch.zip" "$$(basename $$src)"); \
		else \
			tar -czf "dist/$(BINARY)-$$os-$$arch.tar.gz" -C dist "$$(basename $$src)"; \
		fi; \
	done
	@echo "✓ Archives in dist/"

# ── Install (platform-aware) ──────────────────────────────────────
.PHONY: install
install: release
ifeq ($(OS),linux)
	install -Dm755 bin/$(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	@echo "✓ Installed to $(DESTDIR)$(BINDIR)/$(BINARY)"
else ifeq ($(OS),darwin)
	install -d $(DESTDIR)$(BINDIR)
	install -m755 bin/$(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	@echo "✓ Installed to $(DESTDIR)$(BINDIR)/$(BINARY)"
else
	@echo "Run 'make install PREFIX=/c/Users/$$USER/bin' on Windows (Git Bash)"
	@mkdir -p $(DESTDIR)$(BINDIR)
	@cp bin/$(BINARY).exe $(DESTDIR)$(BINDIR)/$(BINARY).exe 2>/dev/null || \
		cp bin/$(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	@echo "✓ Installed to $(DESTDIR)$(BINDIR)/"
endif

# ── Uninstall ──────────────────────────────────────────────────────
.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY)
	@echo "✓ Removed $(DESTDIR)$(BINDIR)/$(BINARY)"

# ── System-wide install (needs sudo/root) ──────────────────────────
.PHONY: install-system
install-system: release
	install -Dm755 bin/$(BINARY) $(DESTDIR)/usr/local/bin/$(BINARY)
	@echo "✓ Installed to $(DESTDIR)/usr/local/bin/$(BINARY)"

# ── Development ────────────────────────────────────────────────────
.PHONY: run
run:
	go run $(MODULE) $(FILE)

.PHONY: test
test:
	$(CGO) go test -count=1 ./...

.PHONY: test-cover
test-cover:
	$(CGO) go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ coverage.html"

.PHONY: lint
lint:
	go vet ./...

.PHONY: fmt
fmt:
	gofmt -w .
	goimports -w . 2>/dev/null || true

# ── Cleanup ────────────────────────────────────────────────────────
.PHONY: clean
clean:
	rm -rf bin/ dist/ coverage.out coverage.html
	@echo "✓ Clean"

# ── Help ───────────────────────────────────────────────────────────
.PHONY: help
help:
	@printf 'Usage: make [target]\n'
	@printf '\n'
	@printf 'Build:\n'
	@printf '  build          Build for current platform\n'
	@printf '  release        Build with version info\n'
	@printf '  build-all      Cross-compile linux/darwin/windows (amd64+arm64)\n'
	@printf '  dist           build-all + tarballs/zips\n'
	@printf '\n'
	@printf 'Install (PREFIX=%s):\n' "$(PREFIX)"
	@printf '  install        Install to $$PREFIX/bin (default: ~/.local/bin)\n'
	@printf '  install-system Install to /usr/local/bin\n'
	@printf '  uninstall      Remove from $$PREFIX/bin\n'
	@printf '\n'
	@printf 'Dev:\n'
	@printf '  run            Run from source (FILE=readme.md)\n'
	@printf '  test           Run tests\n'
	@printf '  test-cover     Tests + coverage report\n'
	@printf '  lint           go vet\n'
	@printf '  fmt            Format code\n'
	@printf '  clean          Remove build artifacts\n'
	@printf '\n'
	@printf 'Variables: VERSION, PREFIX, DESTDIR, FILE\n'
