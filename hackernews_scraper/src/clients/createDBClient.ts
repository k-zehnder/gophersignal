import { connectToDatabase } from '../database/connection';
import { Config } from '../types/index';

// Connects to the MySQL database and returns the db client object
export const createDBClient = async (config: Config) => {
  const db = await connectToDatabase(config.mysql);
  return db;
};
