package nav_test

import (
	"strings"
	"testing"

	"github.com/Lance52259/doc-draft/internal/nav"
)

func TestExcerptAntiDDoSStyle(t *testing.T) {
	index := `# Introduction

## What is Anti-DDoS

Anti-DDoS (Anti-Distributed Denial of Service) is a distributed denial-of-service attack protection service provided by Huawei Cloud, which can effectively protect public IPs from DDoS attacks and ensure stable business operations. Anti-DDoS service provides two protection modes: Basic Protection and Professional Protection. Basic Protection provides free DDoS attack protection capabilities for Huawei Cloud users. When a DDoS attack is detected, the system will automatically start traffic cleaning, filter out attack traffic, and only forward normal traffic to the origin server.

Anti-DDoS service supports protection against multiple attack types.

## Best Practices Overview
`
	want := "Anti-DDoS (Anti-Distributed Denial of Service) is a distributed denial-of-service attack protection service provided by Huawei Cloud, which can effectively protect public IPs from DDoS attacks and ensure stable business operations."
	got := nav.ExcerptWhatIsParagraph(index, nav.EnUS)
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	heading, blurb := nav.ReadmeNavFromIndex(index, nav.EnUS, "ANTI-DDOS")
	if heading != "Anti-DDoS Best Practices" {
		t.Fatalf("heading=%q", heading)
	}
	if blurb != want {
		t.Fatalf("blurb=%q", blurb)
	}
}

func TestExcerptAOMFullFirstParagraph(t *testing.T) {
	index := `# Introduction

## What is Application Operations Management (AOM)

Application Operations Management (AOM) is a one-stop application operations management platform provided by Huawei Cloud, offering enterprises a unified application operations management entry point. AOM helps enterprises achieve automated operations and intelligent management through application monitoring, log management, alarm management, and other functions, improving operational efficiency and quality.

AOM service supports multiple monitoring metrics.

## Best Practices Overview
`
	got := nav.ExcerptWhatIsParagraph(index, nav.EnUS)
	if !strings.Contains(got, "one-stop application operations") || !strings.Contains(got, "operational efficiency and quality.") {
		t.Fatalf("expected full first paragraph, got: %s", got)
	}
	heading, _ := nav.ReadmeNavFromIndex(index, nav.EnUS, "AOM")
	if heading != "Application Operations Management (AOM) Best Practices" {
		t.Fatalf("heading=%q", heading)
	}
}

func TestExcerptChinese(t *testing.T) {
	index := `# 简介

## 什么是Anti-DDoS（Anti-DDoS）

Anti-DDoS（Anti-Distributed Denial of Service）是华为云提供的分布式拒绝服务攻击防护服务，能够有效防护针对公网IP的DDoS攻击，保障业务的稳定运行。Anti-DDoS服务提供基础防护和专业版防护两种防护模式，基础防护为华为云用户提供免费的DDoS攻击防护能力，当检测到DDoS攻击时，系统会自动启动流量清洗，将攻击流量过滤后，仅将正常流量转发给源站服务器。

Anti-DDoS服务支持多种攻击类型的防护。

## 最佳实践简述
`
	got := nav.ExcerptWhatIsParagraph(index, nav.ZhCN)
	if !strings.HasPrefix(got, "Anti-DDoS（Anti-Distributed Denial of Service）是华为云提供的") {
		t.Fatalf("got=%q", got)
	}
	if !strings.Contains(got, "稳定运行。") {
		t.Fatalf("got=%q", got)
	}
	// Second paragraph must not be included
	if strings.Contains(got, "支持多种攻击类型") {
		t.Fatalf("should stop at first paragraph: %q", got)
	}
	heading, _ := nav.ReadmeNavFromIndex(index, nav.ZhCN, "AAD")
	if heading != "Anti-DDoS（Anti-DDoS）最佳实践" {
		t.Fatalf("heading=%q", heading)
	}
}
