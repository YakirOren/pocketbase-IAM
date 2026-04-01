import { defineConfig } from "vitepress";

export default defineConfig({
  title: "pocketbase-IAM",
  description: "AWS IAM-inspired access control for PocketBase",
  base: "/pocketbase-IAM/",
  cleanUrls: true,
  srcExclude: ["internal/**"],
  head: [
    [
      "link",
      {
        rel: "icon",
        type: "image/svg+xml",
        href: "/pocketbase-IAM/favicon.svg",
      },
    ],
  ],
  themeConfig: {
    logo: "/favicon.svg",
    nav: [
      { text: "Guide", link: "/getting-started" },
      { text: "API", link: "/api/setup-options" },
      { text: "Examples", link: "/examples/basic-crud" },
    ],
    sidebar: [
      {
        text: "Get Started",
        items: [
          { text: "What is pocketbase-IAM?", link: "/" },
          { text: "Getting Started", link: "/getting-started" },
        ],
      },
      {
        text: "Concepts",
        items: [
          { text: "Policies", link: "/concepts/policies" },
          { text: "Statements", link: "/concepts/statements" },
          {
            text: "Actions & Resources",
            link: "/concepts/actions-and-resources",
          },
          { text: "Evaluation Flow", link: "/concepts/evaluation-flow" },
          {
            text: "Managed Collections",
            link: "/concepts/managed-collections",
          },
          { text: "Roles", link: "/concepts/roles" },
          { text: "Groups", link: "/concepts/groups" },
          { text: "Caching", link: "/concepts/caching" },
        ],
      },
      {
        text: "Dashboard",
        items: [
          { text: "Overview", link: "/dashboard/overview" },
          { text: "Managing Policies", link: "/dashboard/managing-policies" },
          {
            text: "Managing Roles & Groups",
            link: "/dashboard/managing-roles-and-groups",
          },
          { text: "Policy Simulator", link: "/dashboard/policy-simulator" },
        ],
      },
      {
        text: "API Reference",
        collapsed: true,
        items: [
          { text: "Setup Options", link: "/api/setup-options" },
          { text: "Check Endpoint", link: "/api/check-endpoint" },
          { text: "Simulate Endpoint", link: "/api/simulate-endpoint" },
          { text: "Custom Actions", link: "/api/custom-actions" },
        ],
      },
      {
        text: "Examples",
        collapsed: true,
        items: [
          { text: "Basic CRUD Policy", link: "/examples/basic-crud" },
          { text: "Role-Based Access", link: "/examples/role-based-access" },
          { text: "Group Policies", link: "/examples/group-policies" },
          { text: "Wildcard Patterns", link: "/examples/wildcard-patterns" },
          { text: "Deny Overrides", link: "/examples/deny-overrides" },
        ],
      },
    ],
    socialLinks: [
      {
        icon: "github",
        link: "https://github.com/YakirOren/pocketbase-IAM",
      },
    ],
    search: {
      provider: "local",
    },
  },
});
