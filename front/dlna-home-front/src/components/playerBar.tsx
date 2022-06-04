
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import Box from '@mui/system/Box';
import {
    useContext,
} from 'react';

import CD from '../assets/img/cd.png';
import AppContext from '../store';
export default function PlayerBar() {
    const ctx = useContext(AppContext);
    const store = ctx[0];

    const dispatch = ctx[1];

    const showPlayer = store.player.show;

    if (showPlayer) return null;
    return <>

        <Box sx={{

            height: 56,
            position: 'fixed',
            left: 0,
            right: 0,
            bottom: 56,
            bgcolor: 'background.paper',
            boxShadow: 2,
        }}>
            <Button fullWidth onClick={() => {
                dispatch({
                    type: 'show-player',
                });
            }}>
                <Box sx={{
                    flex: 1,
                    display: 'flex',
                    alignItems: 'center',
                }}>
                    <img src={CD} height="44"/>

                    <Box sx={{
                        pl: 1,
                    }}>
                        <Typography variant='h3' sx={{
                            color: 'text.primary',
                            fontSize: 16,
                            fontWeight: 500,
                        }}>{store.player.currentItem?.name}</Typography>
                    </Box>
                </Box>
            </Button>
        </Box>

    </>;
}
