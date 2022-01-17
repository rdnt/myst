import Vue from "vue";
import VueRouter from "vue-router";
import Sites from "@/views/Sites";
import PasswordGenerator from "@/views/PasswordGenerator";
import Settings from "@/views/Settings";
import Appearance from "@/views/Settings";
import Security from "@/views/Settings";
import Error404 from "@/errors/404";
import EditSite from "@/views/EditSite";
import Entry from "@/components/entry";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home"
    // redirect: "/sites"
  },
  {
    path: "/keystore/:id",
    name: "Keystore",
    component: Sites,
    props: {
      sidebar: true,
      icon:
        "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiPjxwYXRoIGQ9Ik0yMy42MjEgOS4wMTJjLjI0Ny45NTkuMzc5IDEuOTY0LjM3OSAzIDAgNi42MjMtNS4zNzcgMTEuOTg4LTEyIDExLjk4OHMtMTItNS4zNjUtMTItMTEuOTg4YzAtNi42MjMgNS4zNzctMTIgMTItMTIgMi41ODEgMCA0Ljk2OS44MjIgNi45MjcgMi4yMTFsMS43MTgtMi4yMjMgMS45MzUgNi4wMTJoLTYuNThsMS43MDMtMi4yMDRjLTEuNjItMS4xMjgtMy41ODItMS43OTYtNS43MDMtMS43OTYtNS41MiAwLTEwIDQuNDgxLTEwIDEwIDAgNS41MiA0LjQ4IDEwIDEwIDEwIDUuNTE5IDAgMTAtNC40OCAxMC0xMCAwLTEuMDQ1LS4xNjEtMi4wNTMtLjQ1OC0zaDIuMDc5em0tNy42MjEgNy45ODhoLTh2LTZoMXYtMmMwLTEuNjU2IDEuMzQ0LTMgMy0zczMgMS4zNDQgMyAzdjJoMXY2em0tNS04djJoMnYtMmMwLS41NTItLjQ0OC0xLTEtMXMtMSAuNDQ4LTEgMXoiLz48L3N2Zz4="
    },
    children: [
      {
        path: "/keystore/:id/entry/:entryId",
        name: "entry",
        component: Entry,
        props: true
      }
    ]
  },
  {
    path: "/entry/:id(\\d+)/edit",
    name: "EditSite",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: EditSite
  },
  {
    path: "/generator",
    name: "Password Generator",
    component: PasswordGenerator,
    props: {
      sidebar: true,
      icon:
        "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0Ij48cGF0aCBkPSJNMCA2djEyaDI0di0xMmgtMjR6bTIyIDEwaC0yMHYtOGgyMHY4em0tMTQuODEyLTQuMDQ4bC0uMjE3LS42MTktLjY1OC4yNXYtLjcwOGgtLjY0M3YuNzA2bC0uNjQxLS4yNDktLjIxNy42MTkuNjYxLjIyMy0uNDIyLjU1OS41MzcuMzkyLjM5Ny0uNTg2LjQyMi41ODUuNTI3LS4zOTItLjQyMi0uNTU4LjY3Ni0uMjIyem00IDBsLS4yMTctLjYxOS0uNjU4LjI1di0uNzA4aC0uNjQzdi43MDZsLS42NDEtLjI0OS0uMjE3LjYxOS42NjEuMjIzLS40MjIuNTU5LjUzNi4zOTIuMzk4LS41ODYuNDIyLjU4NS41MjctLjM5Mi0uNDIyLS41NTguNjc2LS4yMjJ6bTQgMGwtLjIxNy0uNjE5LS42NTguMjV2LS43MDhoLS42NDN2LjcwNmwtLjY0MS0uMjQ5LS4yMTcuNjE5LjY2MS4yMjMtLjQyMi41NTkuNTM2LjM5Mi4zOTgtLjU4Ni40MjIuNTg1LjUyNy0uMzkyLS40MjItLjU1OC42NzYtLjIyMnptNC44MTIgMi4wNDhoLTN2LTFoM3YxeiIvPjwvc3ZnPg=="
    }
  },
  {
    path: "/settings",
    name: "Settings",
    component: Settings,
    props: {
      sidebar: true,
      icon:
        "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0Ij48cGF0aCBkPSJNMjQgMTQuMTg3di00LjM3NGMtMi4xNDgtLjc2Ni0yLjcyNi0uODAyLTMuMDI3LTEuNTI5LS4zMDMtLjcyOS4wODMtMS4xNjkgMS4wNTktMy4yMjNsLTMuMDkzLTMuMDkzYy0yLjAyNi45NjMtMi40ODggMS4zNjQtMy4yMjQgMS4wNTktLjcyNy0uMzAyLS43NjgtLjg4OS0xLjUyNy0zLjAyN2gtNC4zNzVjLS43NjQgMi4xNDQtLjggMi43MjUtMS41MjkgMy4wMjctLjc1Mi4zMTMtMS4yMDMtLjEtMy4yMjMtMS4wNTlsLTMuMDkzIDMuMDkzYy45NzcgMi4wNTUgMS4zNjIgMi40OTMgMS4wNTkgMy4yMjQtLjMwMi43MjctLjg4MS43NjQtMy4wMjcgMS41Mjh2NC4zNzVjMi4xMzkuNzYgMi43MjUuOCAzLjAyNyAxLjUyOC4zMDQuNzM0LS4wODEgMS4xNjctMS4wNTkgMy4yMjNsMy4wOTMgMy4wOTNjMS45OTktLjk1IDIuNDctMS4zNzMgMy4yMjMtMS4wNTkuNzI4LjMwMi43NjQuODggMS41MjkgMy4wMjdoNC4zNzRjLjc1OC0yLjEzMS43OTktMi43MjMgMS41MzctMy4wMzEuNzQ1LS4zMDggMS4xODYuMDk5IDMuMjE1IDEuMDYybDMuMDkzLTMuMDkzYy0uOTc1LTIuMDUtMS4zNjItMi40OTItMS4wNTktMy4yMjMuMy0uNzI2Ljg4LS43NjMgMy4wMjctMS41Mjh6bS00Ljg3NS43NjRjLS41NzcgMS4zOTQtLjA2OCAyLjQ1OC40ODggMy41NzhsLTEuMDg0IDEuMDg0Yy0xLjA5My0uNTQzLTIuMTYxLTEuMDc2LTMuNTczLS40OS0xLjM5Ni41ODEtMS43OSAxLjY5My0yLjE4OCAyLjg3N2gtMS41MzRjLS4zOTgtMS4xODUtLjc5MS0yLjI5Ny0yLjE4My0yLjg3NS0xLjQxOS0uNTg4LTIuNTA3LS4wNDUtMy41NzkuNDg4bC0xLjA4My0xLjA4NGMuNTU3LTEuMTE4IDEuMDY2LTIuMTguNDg3LTMuNTgtLjU3OS0xLjM5MS0xLjY5MS0xLjc4NC0yLjg3Ni0yLjE4MnYtMS41MzNjMS4xODUtLjM5OCAyLjI5Ny0uNzkxIDIuODc1LTIuMTg0LjU3OC0xLjM5NC4wNjgtMi40NTktLjQ4OC0zLjU3OWwxLjA4NC0xLjA4NGMxLjA4Mi41MzggMi4xNjIgMS4wNzcgMy41OC40ODggMS4zOTItLjU3NyAxLjc4NS0xLjY5IDIuMTgzLTIuODc1aDEuNTM0Yy4zOTggMS4xODUuNzkyIDIuMjk3IDIuMTg0IDIuODc1IDEuNDE5LjU4OCAyLjUwNi4wNDUgMy41NzktLjQ4OGwxLjA4NCAxLjA4NGMtLjU1NiAxLjEyMS0xLjA2NSAyLjE4Ny0uNDg4IDMuNTguNTc3IDEuMzkxIDEuNjg5IDEuNzg0IDIuODc1IDIuMTgzdjEuNTM0Yy0xLjE4OC4zOTgtMi4zMDIuNzkxLTIuODc3IDIuMTgzem0tNy4xMjUtNS45NTFjMS42NTQgMCAzIDEuMzQ2IDMgM3MtMS4zNDYgMy0zIDMtMy0xLjM0Ni0zLTMgMS4zNDYtMyAzLTN6bTAtMmMtMi43NjIgMC01IDIuMjM4LTUgNXMyLjIzOCA1IDUgNSA1LTIuMjM4IDUtNS0yLjIzOC01LTUtNXoiLz48L3N2Zz4="
    },
    children: [
      {
        name: "General",
        path: "/settings",
        component: Appearance
      },
      {
        name: "Appearance",
        path: "/settings/appearance",
        component: Appearance
      },
      {
        name: "Security",
        path: "/settings/security",
        component: Security
      }
    ]
  },
  {
    path: "/assets/*"
  },
  {
    path: "*",
    component: Error404
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

// console.log(__dirname);
export default router;
