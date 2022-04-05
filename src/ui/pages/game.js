import React from 'react';

import {Paper, Box, Divider} from '@mui/material';
import styled from '@emotion/styled';
// import ReactJson from 'react-json-view'
import { useMatch } from "react-location";

import LnDraw from '../components/lndraw';
//import LnChart from '../../components/lnchart';

//(Box)
// const LnBox = styled.div`
//     display: flex;
//     flex-direction: row;
//     justify-content: space-evenly;
// `
// const LnPaper = styled(Paper)`
//     flex:1;
// `

// const LnNumbers = styled(Paper)`
//     width: 500px;
//     flex:1;
//     overflow-y: scroll;
//     margin-bottom: 30px;
// `

// const LnCharts = styled.div`
//     min-width: 1000px;
// `

// const LnTitle = styled.div`
//     font-size: 1rem;
//     font-weight: bold;
// `

const GamePage = (props) => {

    // const router = useRouter();
    // const {data, error} = useSWR(router.query.game, url => getGame(url),{onSuccess: (data, key, config) => {
    //     console.log('swr onsuccess.');
    // }});

    // console.log('LnGamesDetails: ', error);
    // console.log('LnGamesDetails: ', data);

    // if(error) return (<div>{`Failed to load data for URL:${router.query.game}`}</div>)
    // if(!data) return (<div></div>)

    // console.log('LnGamesDetails Router: ', router);

    //const title = V.titleCase(router.query.game).replace('-', ' ');
    const gameData = useMatch().data;

    // const arr = [];
    // const fillarr = (arr) => {

    //     for (var i = 0; i < 100; i++){
    //         arr.push(i);
    //     }
    // }

    // fillarr(arr);

    return (


        <>
            {/* <LnNumbers>
                {
                    arr.map(a => {

                        return (<div key={a}>{`HELLO ${a}`}</div>)
                    
                    })
                }
            </LnNumbers> */}
            {/* <LnNumbers> */}
            <LnDraw data={gameData} />
            {/* </LnNumbers> */}
            {/* <LnCharts>
                <LnTitle>Chart for: {title}</LnTitle>
            </LnCharts> */}
        </>

    )

}

export default GamePage;
