# EasyLPAC
**语言:** [English](./README.md) | [正體中文](./README_zh-TW.md) | [日本語](./README_ja-JP.md)

[lpac](https://github.com/estkme-group/lpac) 图形界面前端

下载: [GitHub Release](https://github.com/creamlike1024/EasyLPAC/releases/latest)

Arch Linux: ![AUR package](https://img.shields.io/aur/version/easylpac) [AUR - easylpac](https://aur.archlinux.org/packages/easylpac)
 感谢 [1ridic](https://github.com/1ridic), [root-core](https://github.com/Root-Core)

NixOS: [NUR](https://github.com/nix-community/NUR#readme) 软件包 https://github.com/nix-community/nur-combined/blob/master/repos/linyinfeng/pkgs/easylpac/default.nix

openSUSE: https://software.opensuse.org/package/easylpac ([OBS](https://build.opensuse.org/package/show/home:Psheng/EasyLPAC))

系统需求:
- Windows 10 及以上（最后一个支持 Windows 7 的版本是 [0.7.7.2](https://github.com/creamlike1024/EasyLPAC/releases/tag/0.7.7.2)）
- 最新版 macOS
- Linux: `pcscd`/`pcsclite`、`libcurl`（供 lpac 使用）和 `gtk3dialog`（供 EasyLPAC 使用）。

目前仅支持 pcsc 的 APDUINTERFACE 和 curl 的 HTTPINTERFACE。

# 用法

运行前请先连接读卡器。

## Linux

lpac 可执行文件搜索顺序：首先在 EasyLPAC 所在目录中搜索；如果找不到，则使用 `/usr/bin/lpac`。

`EasyLPAC-linux-x86_64-with-lpac.tar.gz` 包含预编译的 lpac 可执行文件。如果无法运行，则需要通过包管理器安装 `lpac`，或者自行[编译 lpac](https://github.com/estkme-group/lpac?tab=readme-ov-file#compile)。

## 自动处理通知
EasyLPAC 默认会处理所有操作产生的通知，并在处理成功后自动删除通知。

你可以前往“设置”标签页，取消勾选“自动处理通知”来关闭该行为。

不过，手动操作通知并不符合 GSMA 规范，因此不建议这样做。

# 截图
<p>
<a href="https://github.com/creamlike1024/EasyLPAC/blob/master/screenshots/chipinfo.png"><img src="https://github.com/creamlike1024/EasyLPAC/blob/master/screenshots/chipinfo.png?raw=true"  height="180px"/></a>
<a href="https://github.com/creamlike1024/EasyLPAC/blob/master/screenshots/notification.png"><img src="https://github.com/creamlike1024/EasyLPAC/blob/master/screenshots/notification.png?raw=true" height="180px"/></a>
<a href="https://github.com/creamlike1024/EasyLPAC/blob/master/screenshots/profile.png"><img src="https://github.com/creamlike1024/EasyLPAC/blob/master/screenshots/profile.png?raw=true" height="180px"/></a>
</p>

# 常见问题

## 使用 5ber 时出现 lpac 错误 `euicc_init`

前往“设置 -> lpac ISD-R AID”，点击 5ber 设置 5ber 的自定义 AID，然后重试。

## macOS `SCardTransmit() failed: 80100016`

如果你使用 macOS Sonoma，可能会遇到这个错误：`SCardTransmit() failed: 80100016`

这是因为 Apple 的 USB CCID 读卡器驱动存在 bug。你可以尝试安装读卡器厂商提供的 macOS 驱动，或者阅读下面的文章来解决：

- [Apple's own CCID driver in Sonoma](https://blog.apdu.fr/posts/2023/11/apple-own-ccid-driver-in-sonoma/)
- [macOS Sonoma bug: SCardControl() returns SCARD_E_NOT_TRANSACTED](https://blog.apdu.fr/posts/2023/09/macos-sonoma-bug-scardcontrol-returns-scard_e_not_transacted/)

## `SCardEstablishContext() failed: 8010001D`

这表示 PCSC 服务没有运行。对于 Linux，该服务通常是 `pcscd`。

在基于 systemd 的发行版上启动 `pcscd`：`sudo systemctl start pcscd`

## `SCardListReaders() failed: 8010002E`

读卡器未连接。

## 其他 `SCard` 错误码

完整的 PCSC 错误码说明请参阅 [pcsc-lite: ErrorCodes](https://pcsclite.apdu.fr/api/group__ErrorCodes.html)
