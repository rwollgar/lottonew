import React, { useState } from 'react'


import { css } from '@pigment-css/react';
import {Box} from '@mui/material';

import { If, Choose, Otherwise, When } from "vite-plugin-react-control-statements";
import reactLogo from './assets/react.svg'
import viteLogo from './assets/vite.svg'
import lottoballs_1 from './assets/img/lottoballs_1.jpg'
import './App.css'

function App() {
	
  	const [count, setCount] = useState(0)

	return (
		<>
			<div>
				<a href="https://vitejs.dev" target="_blank">
					<img src={viteLogo} className="logo" alt="Vite logo" />
				</a>
				<a href="https://react.dev" target="_blank">
					<img src={reactLogo} className="logo react" alt="React logo" />
				</a>
				<br/>
				<img src={lottoballs_1} alt="Lotto Balls" />
			</div>
			<Box className={css({backgroundColor:'lightblue',fontWeight:'bold',fontSize:25})}> 
				<Choose>
					<When condition={count < 5}>
						<div className={css({color:'red'})}>Vite + React</div>
					</When>
					<Otherwise>
						<div>{`Vite + React Version: ${React.version}`}</div>
					</Otherwise>
				</Choose>
			</Box>
			<div className="card">
				<button onClick={() => setCount((count) => count + 1)}>
					count is {count}
				</button>
				<p>
					Edit <code>src/App.jsx</code> and save to test HMR
				</p>
			</div>
			<p className="read-the-docs">
				Click on the Vite and React logos to learn more
			</p>
		</>
	)
}

export default App
