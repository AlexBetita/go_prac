import { BotEntry } from "@/lib/types/botTypes";

export const mockBotEntries: BotEntry[] = [
  {
    userMessage: "Are their any posts about programming?",
    type: "related_posts",
    response:
      "Hey there! 🌟 I just dug up some related blog posts, and I think I found a few gems for you. Let's dive in! I've split them into posts that are relevant to programming (yay!) and posts that are... well, not so relevant (still cool but maybe not what we're looking for). Ready? Let's go! 🚀\n\n---\n\n### **Posts Relevant to Programming** 📚\nHere are the golden nuggets that tie back to programming:\n\n| **Slug**             | **Topic**                                      | **Views** | **Created By** |\n|----------------------|-----------------------------------------------|-----------|----------------|\n| `dark-vs-light`      | Dark Mode Devs vs Light Mode Rebels           | 0         | Seeder         |\n| `99-script-problems` | 99 Problems But a Script Ain't One            | 0         | Seeder         |\n| `dev-wrist-health`   | A Dev’s Guide to Healthy Wrists               | 0         | Seeder         |\n\n---\n\n### **Posts Not Relevant to Programming** 🙃 (Still worth a look!)\nThese might not scream \"coding,\" but hey, we're multifaceted humans, right?\n\n| **Slug**              | **Topic**                                  | **Views** | **Created By** |\n|-----------------------|-------------------------------------------|-----------|----------------|\n| `remote-dev-travel`   | Traveling as a Remote Dev: Dream or Distraction? | 0       | Seeder         |\n| `burnout-breakthrough`| From Burnout to Breakthrough              | 0         | Seeder         |\n\n---\n\nThere you have it! Hope these posts spark some curiosity (or at least give you an excuse to procrastinate in style 😄). Which one are you checking out first? ✨",
  },
  {
    userMessage: "How to center a div horizontally?",
    type: "chat",
    response:
      "Use `margin: 0 auto;` or a flex container with `justify-center`.",
  },
  {
    userMessage: "What's the difference between let and const?",
    type: "chat",
    response: "`const` can't be reassigned; `let` can. Both are block‑scoped.",
  },
  {
    userMessage: "How do I add a route in React Router v6?",
    type: "chat",
    response:
      'Wrap your app in `<Routes>` and then `<Route path="/foo" element={<Foo/>}/>`.',
  },
  {
    userMessage: "Explain React’s useEffect hook.",
    type: "chat",
    response:
      "Runs side‑effects after render. The deps array tells it when to re‑run.",
  },
  {
    userMessage: "Reverse a string in JS?",
    type: "chat",
    response: "```js\nstr.split('').reverse().join('')\n```",
  },
  {
    userMessage: "What is a MongoDB aggregation pipeline?",
    type: "chat",
    response:
      "A series of stages (`$match`, `$group`, etc.) that transform your data.",
  },
  {
    userMessage: "How to create a Redux slice?",
    type: "chat",
    response:
      "Use `createSlice({ name, initialState, reducers })` from `@reduxjs/toolkit`.",
  },
  {
    userMessage: "How to deploy to Netlify?",
    type: "chat",
    response:
      "Connect your Git repo and set the build command in Netlify’s dashboard.",
  },
  {
    userMessage: "CSS Grid vs. Flexbox?",
    type: "chat",
    response: "Grid is 2D (rows + columns); Flexbox is 1D (row *or* column).",
  },
  {
    userMessage: "Async/await vs. Promises?",
    type: "chat",
    response: "`async/await` is just cleaner syntax over native Promises.",
  },
  {
    userMessage: "Best way to fetch data in Next.js?",
    type: "chat",
    response: "Use `getStaticProps` for SSG or `getServerSideProps` for SSR.",
  },
  {
    userMessage: "What is a closure?",
    type: "chat",
    response: "A function that ‘remembers’ variables from its lexical scope.",
  },
  {
    userMessage: "Deep clone an object?",
    type: "chat",
    response:
      "Use `structuredClone(obj)` or `JSON.parse(JSON.stringify(obj))`.",
  },
  {
    userMessage: "What’s a CSS pseudo‑class?",
    type: "chat",
    response: "Selectors like `:hover`, `:focus`, targeting element states.",
  },
  {
    userMessage: "Center vertically in CSS?",
    type: "chat",
    response: "Flex container + `items-center` or `translateY(-50%)` trick.",
  },
  {
    userMessage: "REST vs. GraphQL?",
    type: "chat",
    response:
      "REST: fixed endpoints. GraphQL: flexible queries on one endpoint.",
  },
  {
    userMessage: "What is CORS?",
    type: "chat",
    response:
      "Cross‑Origin Resource Sharing: browser policy for external requests.",
  },
  {
    userMessage: "Handle errors in Express?",
    type: "chat",
    response:
      "Use an error middleware: `app.use((err, req, res, next) => { … })`.",
  },
  {
    userMessage: "Write Jest unit tests?",
    type: "chat",
    response:
      "Use `describe()`, `test()` (or `it()`), and `expect()` assertions.",
  },
  {
    userMessage: "Explain OAuth2 flow.",
    type: "chat",
    response: "Auth code → exchange for token → use token for API calls.",
  },
  {
    userMessage: "What’s Tailwind CSS?",
    type: "chat",
    response: "A utility‑first CSS framework with atomic class names.",
  },
  {
    userMessage: "Var vs. let?",
    type: "chat",
    response: "`var` is function‑scoped & hoisted; `let` is block‑scoped.",
  },
  {
    userMessage: "Env vars in React?",
    type: "chat",
    response:
      "Prefix with `REACT_APP_` and access via `process.env.REACT_APP_*`.",
  },
  {
    userMessage: "What is debouncing?",
    type: "chat",
    response: "Delaying a fn until a pause in events (e.g. keystrokes).",
  },
  {
    userMessage: "Explain CSS specificity.",
    type: "chat",
    response: "Inline > ID > class/attr > element selectors.",
  },
  {
    userMessage: "Optimize React renders?",
    type: "chat",
    response: "Use `React.memo`, `useCallback` & avoid anonymous props.",
  },
  {
    userMessage: "What is WebSocket?",
    type: "chat",
    response: "Full‑duplex TCP protocol for real‑time comms.",
  },
  {
    userMessage: "Implement dark mode?",
    type: "chat",
    response:
      "Use `prefers-color-scheme` or toggle a `.dark` class on `<html>`.",
  },
  {
    userMessage: "What is Docker?",
    type: "chat",
    response: "Containerization platform to bundle apps + dependencies.",
  },
  {
    userMessage: "Pre‑commit hooks?",
    type: "chat",
    response: "Use Husky or add scripts in `.git/hooks`.",
  },
  {
    userMessage: "How to paginate in MongoDB?",
    type: "chat",
    response: "Use `.skip(page * limit).limit(limit)` on your query.",
  },
  {
    userMessage: "Explain event bubbling.",
    type: "chat",
    response: "Events propagate from deepest child up through ancestors.",
  },
];
