import React from 'react';

//import {Box, Divider} from '@material-ui/core';
import styled from '@emotion/styled';
//import V from 'voca';
import { useMatch } from "react-location";

import LnDraw from '../components/lndraw';
//import LnChart from '../../components/lnchart';

//(Box)
// const LnBox = styled.div`
//     display: flex;
//     flex-direction: row;
//     justify-content: space-evenly;
// `

const LnNumbers = styled.div`
    max-width: 500px;
`

const LnCharts = styled.div`
    min-width: 1000px;
`

const LnTitle = styled.div`
    font-size: 1rem;
    font-weight: bold;
`

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
    const data = useMatch().gamedata;
    const title = 'title 1';



    return (

        <>
            <LnNumbers>
                <LnTitle>{title}</LnTitle>
                <LnDraw details={JSON.stringify(data)}/>
            </LnNumbers>
            <LnCharts>
                <LnTitle>Chart for: {title}</LnTitle>
            </LnCharts>
        </>

    )

}

export default GamePage;
