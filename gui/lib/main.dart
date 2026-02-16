import 'package:flutter/material.dart';

void main() {
  runApp(const MainApp());
}

class MainApp extends StatelessWidget {
  const MainApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
        colorSchemeSeed: const Color.fromRGBO(200, 200, 200, 1.0),
      ),
      darkTheme: ThemeData(brightness: Brightness.dark),
      title: 'WSL Plus',
      home: const Scaffold(body: MainScreen()),
    );
  }
}

class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  final List<String> _logs = [];

  void _addLog(String message) {
    setState(() {
      _logs.add(message);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: <Widget>[
        Expanded(
          child: Padding(
            padding: const EdgeInsets.all(4.0),
            child: Row(
              children: <Widget>[
                const Spacer(),
                ButtonGroupLeft(enabled: true, onLog: _addLog),
                ButtonGroupRight(enabled: true, onLog: _addLog),
                const Spacer(),
              ],
            ),
          ),
        ),
        Container(
          height: 200,
          decoration: BoxDecoration(
            color: Theme.of(context).colorScheme.surfaceContainerHighest,
            border: Border(
              top: BorderSide(
                color: Theme.of(context).colorScheme.outline,
                width: 1,
              ),
            ),
          ),
          // NOTE: show a log of recent actions
          // TODO: this should be toggled by the user
          child: ListView.builder(
            padding: const EdgeInsets.all(8.0),
            itemCount: _logs.length,
            itemBuilder: (context, index) {
              final opacity = 1.0 - (index * 0.1).clamp(0.0, 0.7);
              return Padding(
                padding: const EdgeInsets.symmetric(vertical: 2.0),
                child: Opacity(
                  opacity: opacity,
                  child: Text(
                    _logs[_logs.length - 1 - index],
                    style: Theme.of(
                      context,
                    ).textTheme.bodyMedium?.copyWith(fontFamily: 'monospace'),
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}

class ButtonGroupLeft extends StatelessWidget {
  const ButtonGroupLeft({
    super.key,
    required this.enabled,
    required this.onLog,
  });

  final bool enabled;
  final Function(String) onLog;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(4.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: <Widget>[
          Column(
            mainAxisSize: MainAxisSize.min,
            children: <Widget>[
              Text(
                'Data',
                style: Theme.of(
                  context,
                ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
              ),
              Container(
                width: 40,
                height: 3,
                decoration: BoxDecoration(
                  color: Theme.of(context).colorScheme.primary,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ],
          ),
          ElevatedButton(
            onPressed: enabled ? () => onLog('listed WSL instances') : null,
            child: const Text('List'),
          ),
          ElevatedButton(
            onPressed: enabled ? () => onLog('info about WSL instances') : null,
            child: const Text('Info'),
          ),
          ElevatedButton(
            onPressed: enabled ? () => onLog('default WSL version') : null,
            child: const Text('Default'),
          ),
        ],
      ),
    );
  }
}

class ButtonGroupRight extends StatelessWidget {
  const ButtonGroupRight({
    super.key,
    required this.enabled,
    required this.onLog,
  });

  final bool enabled;
  final Function(String) onLog;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(4.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: <Widget>[
          Column(
            mainAxisSize: MainAxisSize.min,
            children: <Widget>[
              Text(
                'Actions',
                style: Theme.of(
                  context,
                ).textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
              ),
              Container(
                width: 40,
                height: 3,
                decoration: BoxDecoration(
                  color: Theme.of(context).colorScheme.primary,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ],
          ),
          ElevatedButton(
            onPressed: enabled ? () => onLog('installed WSL instances') : null,
            child: const Text('Install'),
          ),
          ElevatedButton(
            onPressed: enabled ? () => onLog('renamed WSL instances') : null,
            child: const Text('Rename'),
          ),
          ElevatedButton(
            onPressed: enabled ? () => onLog('backed up WSL instances') : null,
            child: const Text('Backup'),
          ),
        ],
      ),
    );
  }
}
