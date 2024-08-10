import { createLazyFileRoute } from '@tanstack/react-router'
import App from '../App'

const Index = () => {
    return (
        <div className="p-2">
            <App/>
        </div>
    )
}

export const Route = createLazyFileRoute('/')({
    component: Index,
})

