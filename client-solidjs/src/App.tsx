import type { Component } from 'solid-js';

import logo from '@/logo.svg';
import styles from '@/App.module.css';

const App: Component = () => {
  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <img src={logo} class={styles.logo} alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <p>HIIIII SOLID</p>
        <a href="https://github.com/login/oauth/authorize?client_id=ffe7776b7ca3b5629066">Login with Github </a>
        <a
          class={styles.link}
          href="https://github.com/solidjs/solid"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn Solid
        </a>
      </header>
    </div>
  );
};

export default App;
