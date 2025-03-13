export type Db = {
    uuid?: string;
    name: string;
    created_at?: string;
    version?: string;
    num_tables?: number;
    file_size?: number;
}

export type Table = {
    name: string;
}