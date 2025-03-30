import {createContext, ReactNode, useContext} from "react";
import {Configuration, TopicsApi} from "../backend_api";

const ApiClientContext = createContext<TopicsApi | undefined>(undefined);

export const useApiClientContext = () => {
    const context = useContext(ApiClientContext);
    if (!context) {
        throw new Error('useApiClientContext must be used within a ApiClientProvider');
    }
    return context;
};

export const ApiClientProvider = ({ children }: { children: ReactNode }) => {
    const apiConfig = new Configuration({
        basePath: '/api'
    })
    const apiClient = new TopicsApi(apiConfig);

    return (
        <ApiClientContext.Provider value={ apiClient }>
            {children}
        </ApiClientContext.Provider>
    );
};