# daemon

Compatibility note: When cross-compiling with TinyGo for embedded or router-class Linux targets, this package is generally safe for Linux userland targets (amd64, arm64, mips*, riscv) if the target provides a standard libc and shell environment. Caveats:

- If a package uses os/user, syscall ioctl, or direct filesystem assumptions, it may fail to compile or behave differently under TinyGo. Test builds on your TinyGo target before deploying.
- For maximal TinyGo portability, prefer environment-variable fallbacks and avoid os/user or heavy syscall usage.

If you need help making this package TinyGo-friendly, I can add build-tagbed fallbacks.
