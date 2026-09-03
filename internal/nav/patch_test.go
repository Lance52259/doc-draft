package nav_test

import (
	"strings"
	"testing"

	"github.com/Lance52259/doc-draft/internal/nav"
)

const sampleSUMMARY = `# Summary

* [产品介绍](introductions/)
  * [Terraform简介](introductions/terraform_introduction.md)
* [最佳实践](best-practices/)
  * [简介](best-practices/README.md)
  * [Anti-DDoS](best-practices/anti-ddos/)
    * [简介](best-practices/anti-ddos/index.md)
    * [部署基础防护](best-practices/anti-ddos/basic.md)
  * [AOM](best-practices/aom/)
    * [简介](best-practices/aom/index.md)
    * [部署AOM告警动作回调](best-practices/aom/action_callback.md)
* [贡献](contribute.md)
`

func TestPatchSUMMARYExistingService(t *testing.T) {
	got, err := nav.PatchSUMMARY(sampleSUMMARY, "aom", "AOM", "prevent_elb_alarm_storm", "部署AOM防止ELB告警风暴")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "best-practices/aom/prevent_elb_alarm_storm.md") {
		t.Fatalf("missing practice:\n%s", got)
	}
	if !strings.Contains(got, "best-practices/anti-ddos/basic.md") || !strings.Contains(got, "* [贡献]") {
		t.Fatalf("destroyed existing entries:\n%s", got)
	}
	if strings.Count(got, "# Summary") != 1 {
		t.Fatalf("title broken:\n%s", got)
	}
}

func TestPatchSUMMARYNewService(t *testing.T) {
	got, err := nav.PatchSUMMARY(sampleSUMMARY, "aad", "AAD", "black_white_lists", "部署黑白名单防护")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "  * [AAD](best-practices/aad/)") {
		t.Fatalf("missing service:\n%s", got)
	}
	// aad should appear before anti-ddos alphabetically
	aad := strings.Index(got, "best-practices/aad/")
	anti := strings.Index(got, "best-practices/anti-ddos/")
	if aad < 0 || anti < 0 || aad > anti {
		t.Fatalf("service order wrong:\n%s", got)
	}
	if !strings.Contains(got, "best-practices/aom/action_callback.md") {
		t.Fatalf("lost aom:\n%s", got)
	}
}

func TestPatchServiceIndex(t *testing.T) {
	base := `# 简介

## 最佳实践列表

本章节包含以下最佳实践：

* [部署基础防护](basic.md) - 介绍基础防护。
* [部署LTS配置](lts_config.md) - 介绍LTS。

## 参考资料
`
	got, err := nav.PatchServiceIndex(base, "default_protection_policy", "部署默认防护策略", "介绍默认防护策略")
	if err != nil {
		t.Fatal(err)
	}
	basic := strings.Index(got, "(basic.md)")
	def := strings.Index(got, "(default_protection_policy.md)")
	lts := strings.Index(got, "(lts_config.md)")
	if basic < 0 || def < 0 || lts < 0 || !(basic < def && def < lts) {
		t.Fatalf("order:\n%s", got)
	}
}

func TestPatchBestPracticesREADME(t *testing.T) {
	base := `# 中心

## 文档导航

### [Anti-DDoS最佳实践](anti-ddos/index.md)

Anti-DDoS 简介。

### [应用运维管理（AOM）最佳实践](aom/index.md)

AOM 简介。

## 其他
`
	got, err := nav.PatchBestPracticesREADME(base, "aad", "DDoS高防（AAD）最佳实践", "AAD 简介。")
	if err != nil {
		t.Fatal(err)
	}
	aad := strings.Index(got, "(aad/index.md)")
	anti := strings.Index(got, "(anti-ddos/index.md)")
	if aad < 0 || anti < 0 || aad > anti {
		t.Fatalf("nav order:\n%s", got)
	}
	if !strings.Contains(got, "## 其他") {
		t.Fatalf("lost suffix:\n%s", got)
	}
}

func TestIsDestructiveUpdate(t *testing.T) {
	if !nav.IsDestructiveUpdate(sampleSUMMARY, "# 目录\n\n* [AAD](best-practices/aad/)\n") {
		t.Fatal("expected destructive")
	}
	patched, _ := nav.PatchSUMMARY(sampleSUMMARY, "aad", "AAD", "black_white_lists", "部署黑白名单防护")
	if nav.IsDestructiveUpdate(sampleSUMMARY, patched) {
		t.Fatal("surgical patch should not be destructive")
	}
}
