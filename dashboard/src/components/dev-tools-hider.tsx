"use client"

import { useState } from "react"

export function DevToolsHider() {
  // TODO: Replace this hardcoded state with your actual user role context/store
  const [userRole] = useState("admin") // Change to "superadmin" to see the Dev Tools badge

  // Only run logic in development mode
  if (process.env.NODE_ENV !== "development") {
    return null
  }

  // If the user is a superadmin, don't hide the dev tools
  if (userRole === "superadmin") {
    return null
  }

  // For all other roles, inject CSS to hide the dev tools badge
  return (
    <style
      dangerouslySetInnerHTML={{
        __html: `
          .__web-inspector-hide-shortcut__,
          [data-nextjs-dev-tools-button="true"],
          [data-next-badge="true"] {
            display: none !important;
          }
        `,
      }}
    />
  )
}
