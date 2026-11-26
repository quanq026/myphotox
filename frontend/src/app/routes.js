/*

Copyright (c) 2018 - 2026 PhotoPrism UG. All rights reserved.

    This program is free software: you can redistribute it and/or modify
    it under Version 3 of the GNU Affero General Public License (the "AGPL"):
    <https://docs.photoprism.app/license/agpl>

*/

import Photos from "page/photos.vue";
import Albums from "page/albums.vue";
import AlbumPhotos from "page/album/photos.vue";
import Settings from "page/settings.vue";
import Login from "page/auth/login.vue";
import About from "page/about/about.vue";
import License from "page/about/license.vue";
import { $gettext, $pgettext } from "common/gettext";
import { $config, $session } from "./session";

const c = window.__CONFIG__;
const siteTitle = c.siteTitle ? c.siteTitle : c.name;
const loginRoute = "login";

export function safeReturnTo(value) {
  if (!value || typeof value !== "string") {
    return "";
  }
  const trimmed = value.trim();
  const here = typeof window !== "undefined" ? window.location?.origin : "";
  if (!trimmed || !here) {
    return "";
  }
  try {
    const parsed = new URL(trimmed, here);
    if (parsed.origin === here && !parsed.pathname.startsWith("//")) {
      return parsed.pathname + parsed.search + parsed.hash;
    }
  } catch {
    // Fall through to reject.
  }
  return "";
}

export default [
  {
    name: "home",
    path: "/",
    redirect: () => {
      return { name: $session.getDefaultRoute() };
    },
  },
  {
    name: "about",
    path: "/about",
    component: About,
    meta: { title: $gettext("About"), requiresAuth: false },
  },
  {
    name: "license",
    path: "/license",
    component: License,
    meta: { title: $gettext("License"), requiresAuth: false },
  },
  {
    name: loginRoute,
    path: "/login",
    component: Login,
    meta: { title: siteTitle, requiresAuth: false, hideNav: true },
    beforeEnter: (to, from, next) => {
      const returnTo = safeReturnTo(to.query?.return_to);
      if (returnTo) {
        $session.setLoginRedirectUrl(returnTo);
      }

      if ($session.loginRequired()) {
        next();
        return;
      }

      if ($session.hasLoginRedirectUrl()) {
        if ($session.loginRedirectLooping()) {
          $session.clearLoginRedirectAttempt();
          $session.clearLoginRedirectUrl();
          next({ name: $session.getDefaultRoute() });
        } else {
          $session.markLoginRedirectAttempt();
          $session.followLoginRedirectUrl();
          next(false);
        }
      } else {
        next({ name: $session.getDefaultRoute() });
      }
    },
  },
  {
    name: "logout",
    path: "/logout",
    meta: { title: siteTitle, requiresAuth: false, hideNav: true },
    beforeEnter: (to, from, next) => {
      $session.signOut();
      next({ name: loginRoute });
    },
  },
  {
    name: "browse",
    path: "/browse",
    component: Photos,
    meta: { title: $gettext("Pictures"), icon: true, requiresAuth: true },
    props: { staticFilter: { photo: "true" } },
    beforeEnter: (to, from, next) => {
      if ($session.loginRequired()) {
        next({ name: loginRoute });
      } else if ($config.deny("photos", "search")) {
        next({ name: $session.getDefaultRoute() });
      } else {
        next();
      }
    },
  },
  {
    name: "all",
    path: "/all",
    component: Photos,
    meta: { title: siteTitle, requiresAuth: true },
    props: { staticFilter: { quality: "0" } },
  },
  {
    name: "photos",
    path: "/photos",
    component: Photos,
    meta: { title: $gettext("Photos"), requiresAuth: true },
    props: { staticFilter: { photo: "true" } },
  },
  {
    name: "albums",
    path: "/albums",
    component: Albums,
    meta: { title: $gettext("Albums"), requiresAuth: true },
    props: { view: "album", defaultOrder: "favorites", staticFilter: { type: "album" } },
  },
  {
    name: "album",
    path: "/albums/:album/:slug",
    component: AlbumPhotos,
    meta: { collectionTitle: "Albums", collectionRoute: "albums", requiresAuth: true },
  },
  {
    name: "folders",
    path: "/folders",
    component: Albums,
    meta: { title: $gettext("Folders"), requiresAuth: true },
    props: { view: "folder", defaultOrder: "name", staticFilter: { type: "folder" } },
  },
  {
    name: "folder",
    path: "/folders/:album/:slug",
    component: AlbumPhotos,
    meta: { collectionTitle: "Folders", collectionRoute: "folders", requiresAuth: true },
  },
  {
    name: "favorites",
    path: "/favorites",
    component: Photos,
    meta: { title: $gettext("Favorites"), requiresAuth: true },
    props: { staticFilter: { favorite: "true" } },
  },
  {
    name: "videos",
    path: "/videos",
    component: Photos,
    meta: { title: $gettext("Videos"), requiresAuth: true },
    props: { staticFilter: { video: "true" } },
  },
  {
    name: "archive",
    path: "/archive",
    component: Photos,
    meta: {
      title: $pgettext("Noun", "Archive"),
      requiresAuth: true,
    },
    props: { staticFilter: { archived: "true", public: "" } },
    beforeEnter: (to, from, next) => {
      if ($session.loginRequired()) {
        next({ name: loginRoute });
      } else if ($config.deny("photos", "delete")) {
        next({ name: $session.getDefaultRoute() });
      } else {
        next();
      }
    },
  },
  {
    name: "settings",
    path: "/settings",
    component: Settings,
    meta: {
      title: $gettext("Settings"),
      requiresAuth: true,
      settings: true,
      background: "background",
    },
    props: { tab: "settings_general" },
  },
  {
    path: "/settings/content",
    redirect: "/settings",
  },
  {
    path: "/settings/media",
    redirect: "/settings",
  },
  {
    name: "settings_account",
    path: "/settings/account",
    component: Settings,
    meta: {
      title: $gettext("Settings"),
      requiresAuth: true,
      settings: true,
      background: "background",
    },
    props: { tab: "settings_account" },
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/browse",
  },
];
