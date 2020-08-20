export interface AI {
    /**
     * example:
     * 1
     */
    id: number;
    /**
     * example:
     * cccccc
     */
    commit: string;
    ai_github: AIGithub;
}
export interface AIGithub {
    /**
     * example:
     * 1
     */
    id: number;
    /**
     * example:
     * wass80/reversi-random
     */
    github: string;
    /**
     * example:
     * master
     */
    branch: string;
    user: User;
    game: Game;
}
export interface Game {
    /**
     * example:
     * 1
     */
    id: number;
    /**
     * example:
     * reversi
     */
    name: string;
}
export interface Match {
    /**
     * example:
     * 1
     */
    id: number;
    game: Game;
    /**
     * example:
     * running
     */
    state: string;
    /**
     * example:
     * exception
     */
    exception: string;
    /**
     * example:
     * put 1
     * 
     */
    record?: string;
    results: {
        ai?: AI;
        /**
         * example:
         * 12
         */
        result?: number;
        /**
         * example:
         * exception
         */
        exception?: string;
        /**
         * example:
         * stderr
         */
        stderr?: string;
    }[];
}
export interface User {
    /**
     * example:
     * 1
     */
    id: number;
    /**
     * example:
     * wass80
     */
    name: string;
}
