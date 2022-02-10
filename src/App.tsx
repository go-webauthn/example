import React from 'react';
import './App.css';
import IndexView from "./views/IndexView";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import { orange } from "@mui/material/colors";

function App() {
  return (
      <ThemeProvider theme={theme}>
        <IndexView />
      </ThemeProvider>
  );
}

export default App;

const theme = createTheme({
  status: {
    danger: orange[500],
  }
})