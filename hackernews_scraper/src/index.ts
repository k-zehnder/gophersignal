// Initializes clients, assembles services, and runs the article workflow.

import config from './config/config';
import { createClients } from './clients/createClients';
import { createServices } from './services/createServices';
import { Workflow } from './workflow';

const EXIT_CODE_SUCCESS = 0;
const EXIT_CODE_FAILURE = 1;

const main = async (): Promise<void> => {
  let workflow: Workflow | null = null;
  try {
    // Create low-level infrastructure clients
    const clients = await createClients(config);

    // Assemble high-level services
    const services = createServices(clients);

    // Run the workflow
    workflow = new Workflow(services);
    await workflow.run();

    console.info('Workflow completed successfully');
    process.exit(EXIT_CODE_SUCCESS);
  } catch (error) {
    console.error('Workflow execution error:', error);
    process.exit(EXIT_CODE_FAILURE);
  } finally {
    if (workflow) await workflow.shutdown();
  }
};

if (require.main === module) main();
