/**
 * Assistant Test Panel Plugin
 * Created for branch testing by Grafana Assistant.
 */
export const pluginId = 'assistant-test';
export const pluginVersion = '1.0.0';

export interface AssistantTestOptions {
  title: string;
  enabled: boolean;
}

export const defaults: AssistantTestOptions = {
  title: 'Assistant Test',
  enabled: true,
};
