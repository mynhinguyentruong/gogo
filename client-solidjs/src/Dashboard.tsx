import { useParams, useSearchParams } from "@solidjs/router"
import { create } from "domain"
import { Component, createEffect, createSignal } from "solid-js"

import FingerprintJS from '@fingerprintjs/fingerprintjs';

// async function getUser(code) {
//   // make req to backend server
//   // /api/auth/callback?code=asdkjhasdjk
//   const response = await fetch(`http://localhost:8080/api/auth/callback?code=${code}`)
// }



const Dashboard: Component = () => {
  const [visitorId, setVisitorId] = createSignal("")

  createEffect(() => {
    const location = window.location
    console.log({location})
    const setFp = async () => {
      const fp = await FingerprintJS.load();

      const { visitorId } = await fp.get();

      setVisitorId(visitorId);
    }

    setFp()

  })
  return (
    <div>
      <h1>Visitor Id: {visitorId()}</h1>
      <p>Last visit: </p>
      <p>You have visited this page x times</p>
    </div>
  )
}

export default Dashboard
