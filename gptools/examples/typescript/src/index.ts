import { dag, object, func, Secret, Directory } from "@dagger.io/dagger";

@object()
class Examples {
  /**
   * example on how to run a full e2e RAG model across different types of
   * documents in a directory
   */
  @func()
  gptools_rag(openaiApiKey: Secret, question: string): Promise<string> {
    const nixPaper = dag.http(
      "https://edolstra.github.io/pubs/nspfssd-lisa2004-final.pdf",
    );
    const foxImage = dag.http(
      "https://fsquaredmarketing.com/wp-content/uploads/2024/04/bitter-font.png",
    );
    return dag
      .gptools()
      .rag(
        openaiApiKey,
        dag
          .directory()
          .withFile("nix-paper.pdf", nixPaper)
          .withFile("image.png", foxImage),
        question,
      );
  }

  /**
   * example on how to return a transcript from a video file
   */
  @func()
  async gptoolsTranscript(
    // free to use movie https://mango.blender.org/about/
    url = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4",
  ): Promise<string> {
    const video = dag.http(url);
    return dag.gptools().transcript(video).contents();
  }

  /**
   * example on how to get a transcript for a video and
   * return it as a *Directory to use in other functions
   */
  @func()
  gptoolsTranscript_Directory(
    // free to use movie https://mango.blender.org/about/
    url = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4",
  ): Directory {
    const video = dag.http(url);
    return dag
      .directory()
      .withFile("video-transcript.txt", dag.gptools().transcript(video));
  }
}
